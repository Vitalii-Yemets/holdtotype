#include <cctype>
#include <chrono>
#include <cstdio>
#include <cstdlib>
#include <iostream>
#include <map>
#include <memory>
#include <sstream>
#include <string>
#include <utility>
#include <vector>

#include "translator/parser.h"
#include "translator/response.h"
#include "translator/response_options.h"
#include "translator/service.h"
#include "translator/translation_model.h"

using namespace marian::bergamot;
using Clock = std::chrono::steady_clock;

namespace {

struct Value {
  enum Kind { Null, Bool, Number, String, Array, Object } kind{Null};
  bool b{false};
  double n{0};
  std::string s;
  std::vector<Value> items;
  std::vector<std::pair<std::string, Value>> fields;

  const Value *get(const std::string &key) const {
    for (const auto &f : fields) {
      if (f.first == key) return &f.second;
    }
    return nullptr;
  }
  std::string str(const std::string &key) const {
    const Value *v = get(key);
    return v && v->kind == String ? v->s : std::string();
  }
  bool flag(const std::string &key) const {
    const Value *v = get(key);
    return v && v->kind == Bool && v->b;
  }
  double num(const std::string &key) const {
    const Value *v = get(key);
    return v && v->kind == Number ? v->n : 0;
  }
};

class Reader {
 public:
  explicit Reader(const std::string &text) : t_(text) {}

  bool parse(Value &out) {
    space();
    if (!value(out)) return false;
    space();
    return i_ == t_.size();
  }

 private:
  const std::string &t_;
  size_t i_{0};

  void space() {
    while (i_ < t_.size() && (t_[i_] == ' ' || t_[i_] == '\t' || t_[i_] == '\r' || t_[i_] == '\n')) ++i_;
  }

  bool literal(const char *word) {
    size_t n = std::char_traits<char>::length(word);
    if (t_.compare(i_, n, word) != 0) return false;
    i_ += n;
    return true;
  }

  static void put(std::string &s, unsigned cp) {
    if (cp < 0x80) {
      s += static_cast<char>(cp);
    } else if (cp < 0x800) {
      s += static_cast<char>(0xC0 | (cp >> 6));
      s += static_cast<char>(0x80 | (cp & 0x3F));
    } else if (cp < 0x10000) {
      s += static_cast<char>(0xE0 | (cp >> 12));
      s += static_cast<char>(0x80 | ((cp >> 6) & 0x3F));
      s += static_cast<char>(0x80 | (cp & 0x3F));
    } else {
      s += static_cast<char>(0xF0 | (cp >> 18));
      s += static_cast<char>(0x80 | ((cp >> 12) & 0x3F));
      s += static_cast<char>(0x80 | ((cp >> 6) & 0x3F));
      s += static_cast<char>(0x80 | (cp & 0x3F));
    }
  }

  bool hex4(unsigned &cp) {
    if (i_ + 4 > t_.size()) return false;
    cp = 0;
    for (int k = 0; k < 4; ++k) {
      char c = t_[i_++];
      cp <<= 4;
      if (c >= '0' && c <= '9') cp |= c - '0';
      else if (c >= 'a' && c <= 'f') cp |= c - 'a' + 10;
      else if (c >= 'A' && c <= 'F') cp |= c - 'A' + 10;
      else return false;
    }
    return true;
  }

  bool string(std::string &s) {
    if (i_ >= t_.size() || t_[i_] != '"') return false;
    ++i_;
    while (i_ < t_.size()) {
      char c = t_[i_++];
      if (c == '"') return true;
      if (c != '\\') {
        s += c;
        continue;
      }
      if (i_ >= t_.size()) return false;
      char e = t_[i_++];
      switch (e) {
        case '"': s += '"'; break;
        case '\\': s += '\\'; break;
        case '/': s += '/'; break;
        case 'b': s += '\b'; break;
        case 'f': s += '\f'; break;
        case 'n': s += '\n'; break;
        case 'r': s += '\r'; break;
        case 't': s += '\t'; break;
        case 'u': {
          unsigned cp;
          if (!hex4(cp)) return false;
          if (cp >= 0xD800 && cp < 0xDC00 && i_ + 1 < t_.size() && t_[i_] == '\\' && t_[i_ + 1] == 'u') {
            i_ += 2;
            unsigned lo;
            if (!hex4(lo)) return false;
            cp = 0x10000 + ((cp - 0xD800) << 10) + (lo - 0xDC00);
          }
          put(s, cp);
          break;
        }
        default: return false;
      }
    }
    return false;
  }

  bool value(Value &v) {
    space();
    if (i_ >= t_.size()) return false;
    char c = t_[i_];
    if (c == '"') {
      v.kind = Value::String;
      return string(v.s);
    }
    if (c == '{') {
      v.kind = Value::Object;
      ++i_;
      space();
      if (i_ < t_.size() && t_[i_] == '}') {
        ++i_;
        return true;
      }
      while (true) {
        space();
        std::string key;
        if (!string(key)) return false;
        space();
        if (i_ >= t_.size() || t_[i_] != ':') return false;
        ++i_;
        Value item;
        if (!value(item)) return false;
        v.fields.emplace_back(std::move(key), std::move(item));
        space();
        if (i_ < t_.size() && t_[i_] == ',') {
          ++i_;
          continue;
        }
        if (i_ < t_.size() && t_[i_] == '}') {
          ++i_;
          return true;
        }
        return false;
      }
    }
    if (c == '[') {
      v.kind = Value::Array;
      ++i_;
      space();
      if (i_ < t_.size() && t_[i_] == ']') {
        ++i_;
        return true;
      }
      while (true) {
        Value item;
        if (!value(item)) return false;
        v.items.push_back(std::move(item));
        space();
        if (i_ < t_.size() && t_[i_] == ',') {
          ++i_;
          continue;
        }
        if (i_ < t_.size() && t_[i_] == ']') {
          ++i_;
          return true;
        }
        return false;
      }
    }
    if (literal("true")) {
      v.kind = Value::Bool;
      v.b = true;
      return true;
    }
    if (literal("false")) {
      v.kind = Value::Bool;
      return true;
    }
    if (literal("null")) return true;
    size_t start = i_;
    while (i_ < t_.size() && (std::isdigit(static_cast<unsigned char>(t_[i_])) || t_[i_] == '-' || t_[i_] == '+' ||
                              t_[i_] == '.' || t_[i_] == 'e' || t_[i_] == 'E'))
      ++i_;
    if (start == i_) return false;
    v.kind = Value::Number;
    v.n = std::strtod(t_.substr(start, i_ - start).c_str(), nullptr);
    return true;
  }
};

std::string quote(const std::string &s) {
  std::string out = "\"";
  for (unsigned char c : s) {
    switch (c) {
      case '"': out += "\\\""; break;
      case '\\': out += "\\\\"; break;
      case '\n': out += "\\n"; break;
      case '\r': out += "\\r"; break;
      case '\t': out += "\\t"; break;
      default:
        if (c < 0x20) {
          char buf[8];
          std::snprintf(buf, sizeof buf, "\\u%04x", c);
          out += buf;
        } else {
          out += static_cast<char>(c);
        }
    }
  }
  return out + "\"";
}

std::string yamlPath(const std::string &p) { return quote(p); }

double since(Clock::time_point a) { return std::chrono::duration<double, std::milli>(Clock::now() - a).count(); }

void reply(const std::string &line) {
  std::fwrite(line.data(), 1, line.size(), stdout);
  std::fputc('\n', stdout);
  std::fflush(stdout);
}

std::string idField(const Value &req) {
  const Value *v = req.get("id");
  if (!v || v->kind != Value::Number) return "\"id\":0";
  char buf[32];
  std::snprintf(buf, sizeof buf, "\"id\":%.0f", v->n);
  return buf;
}

void fail(const Value &req, const std::string &why) { reply("{" + idField(req) + ",\"ok\":false,\"error\":" + quote(why) + "}"); }

}  // namespace

int main() {
  std::ios::sync_with_stdio(false);
  BlockingService::Config serviceConfig;
  serviceConfig.cacheSize = 0;
  serviceConfig.logger.level = "off";
  BlockingService service(serviceConfig);
  std::map<std::string, std::shared_ptr<TranslationModel>> models;
  reply("{\"ready\":true}");
  std::string line;
  while (std::getline(std::cin, line)) {
    if (!line.empty() && line.back() == '\r') line.pop_back();
    if (line.empty()) continue;
    Value req;
    Reader reader(line);
    if (!reader.parse(req) || req.kind != Value::Object) {
      reply("{\"id\":0,\"ok\":false,\"error\":\"the request is not a JSON object\"}");
      continue;
    }
    std::string op = req.str("op");
    std::string pair = req.str("from") + "-" + req.str("to");
    auto t0 = Clock::now();
    try {
      if (op == "quit") {
        reply("{" + idField(req) + ",\"ok\":true}");
        break;
      } else if (op == "load") {
        std::ostringstream cfg;
        cfg << "relative-paths: false\n"
            << "models:\n  - " << yamlPath(req.str("model")) << "\n"
            << "vocabs:\n  - " << yamlPath(req.str("vocab")) << "\n  - " << yamlPath(req.str("vocab")) << "\n"
            << "shortlist:\n  - " << yamlPath(req.str("lex")) << "\n  - false\n"
            << "beam-size: 1\nnormalize: 1.0\nword-penalty: 0\nmax-length-break: 128\nmini-batch-words: 1024\n"
            << "workspace: 128\nmax-length-factor: 2.0\nskip-cost: true\ncpu-threads: 0\nquiet: true\n"
            << "quiet-translation: true\ngemm-precision: int8shiftAlphaAll\nalignment: soft\n";
        auto options = parseOptionsFromString(cfg.str(), false);
        auto model = std::make_shared<TranslationModel>(options);
        std::vector<std::string> warm{"Hello."};
        std::vector<ResponseOptions> ro(1);
        service.translateMultiple(model, std::move(warm), ro);
        models[pair] = model;
        char buf[64];
        std::snprintf(buf, sizeof buf, ",\"ok\":true,\"ms\":%.1f}", since(t0));
        reply("{" + idField(req) + buf);
      } else if (op == "unload") {
        models.erase(pair);
        reply("{" + idField(req) + ",\"ok\":true}");
      } else if (op == "translate") {
        std::string from = req.str("from"), to = req.str("to");
        ResponseOptions options;
        options.HTML = req.flag("html");
        options.alignment = options.HTML;
        std::vector<ResponseOptions> ro{options};
        std::vector<std::string> source{req.str("text")};
        std::vector<Response> out;
        auto direct = models.find(pair);
        if (direct != models.end()) {
          out = service.translateMultiple(direct->second, std::move(source), ro);
        } else {
          auto first = models.find(from + "-en");
          auto second = models.find("en-" + to);
          if (first == models.end() || second == models.end()) {
            fail(req, "no model is loaded for " + pair);
            continue;
          }
          out = service.pivotMultiple(first->second, second->second, std::move(source), ro);
        }
        char buf[64];
        std::snprintf(buf, sizeof buf, ",\"ms\":%.1f}", since(t0));
        reply("{" + idField(req) + ",\"ok\":true,\"text\":" + quote(out.empty() ? std::string() : out[0].target.text) + buf);
      } else {
        fail(req, "unknown op " + op);
      }
    } catch (const std::exception &e) {
      fail(req, e.what());
    }
  }
  return 0;
}
