# syntax=docker/dockerfile:1

ARG WHISPER_CPP_VERSION=v1.7.4
ARG WHISPER_MODEL=medium-q5_0

FROM debian:bookworm-slim AS whisper-build
RUN apt-get update && apt-get install -y --no-install-recommends \
        git ca-certificates cmake make \
        gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64 \
    && rm -rf /var/lib/apt/lists/*
ARG WHISPER_CPP_VERSION
RUN git clone --depth 1 --branch ${WHISPER_CPP_VERSION} \
        https://github.com/ggerganov/whisper.cpp /src
COPY build/mingw-toolchain.cmake /toolchain.cmake
RUN cmake -S /src -B /bld \
        -DCMAKE_TOOLCHAIN_FILE=/toolchain.cmake \
        -DCMAKE_BUILD_TYPE=Release \
        -DBUILD_SHARED_LIBS=OFF \
        -DWHISPER_BUILD_TESTS=OFF \
        -DWHISPER_BUILD_EXAMPLES=ON \
        -DGGML_NATIVE=OFF \
        -DGGML_AVX=ON -DGGML_AVX2=ON -DGGML_FMA=ON -DGGML_F16C=ON \
        -DGGML_OPENMP=OFF \
        -DCMAKE_EXE_LINKER_FLAGS="-static" \
    && cmake --build /bld --target whisper-server -j"$(nproc)"

FROM debian:trixie-slim AS llama-src
RUN apt-get update && apt-get install -y --no-install-recommends \
        git ca-certificates cmake make g++ \
        gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64 \
    && rm -rf /var/lib/apt/lists/*
ARG LLAMA_CPP_VERSION=v0.2.0
RUN git clone --depth 1 --branch ${LLAMA_CPP_VERSION} \
        https://github.com/ggml-org/llama.cpp /src
COPY build/mingw-toolchain.cmake /toolchain.cmake

FROM llama-src AS llama-build
RUN cmake -S /src -B /bld \
        -DCMAKE_TOOLCHAIN_FILE=/toolchain.cmake \
        -DCMAKE_BUILD_TYPE=Release \
        -DBUILD_SHARED_LIBS=OFF \
        -DLLAMA_BUILD_TESTS=OFF \
        -DLLAMA_BUILD_EXAMPLES=OFF \
        -DLLAMA_BUILD_SERVER=ON \
        -DLLAMA_CURL=OFF \
        -DLLAMA_BUILD_UI=OFF \
        -DLLAMA_USE_PREBUILT_UI=OFF \
        -DHOST_CXX_COMPILER=/usr/bin/g++ \
        -DGGML_NATIVE=OFF \
        -DGGML_AVX=ON -DGGML_AVX2=ON -DGGML_FMA=ON -DGGML_F16C=ON \
        -DGGML_OPENMP=OFF \
        -DCMAKE_EXE_LINKER_FLAGS="-static" \
    && cmake --build /bld --target llama-server -j"$(nproc)"

FROM llama-src AS llama-gpu-build
RUN apt-get update && apt-get install -y --no-install-recommends \
        glslc python3 spirv-headers \
    && rm -rf /var/lib/apt/lists/*
ARG VULKAN_SDK_TAG=v1.4.309
ARG VULKAN_HEADERS_COMMIT=952f776f6573aafbb62ea717d871cd1d6816c387
ARG VULKAN_LOADER_COMMIT=809941a4ca137df69dc9c6e8eb456bd70309197c
RUN git clone --depth 1 --branch ${VULKAN_SDK_TAG} https://github.com/KhronosGroup/Vulkan-Headers /vkh \
    && test "$(git -C /vkh rev-parse HEAD)" = "${VULKAN_HEADERS_COMMIT}" \
    && git clone --depth 1 --branch ${VULKAN_SDK_TAG} https://github.com/KhronosGroup/Vulkan-Loader /vkl \
    && test "$(git -C /vkl rev-parse HEAD)" = "${VULKAN_LOADER_COMMIT}"
RUN cmake -S /vkh -B /vkh/b -DCMAKE_INSTALL_PREFIX=/vk \
    && cmake --install /vkh/b \
    && mkdir -p /vk/lib \
    && x86_64-w64-mingw32-dlltool -d /vkl/loader/vulkan-1.def -D vulkan-1.dll -l /vk/lib/libvulkan-1.dll.a \
    && cp -r /usr/include/spirv /vk/include/ \
    && mkdir -p /vk/share/cmake \
    && cp -r /usr/share/cmake/SPIRV-Headers /vk/share/cmake/
RUN cmake -S /src -B /bld \
        -DCMAKE_TOOLCHAIN_FILE=/toolchain.cmake \
        -DCMAKE_BUILD_TYPE=Release \
        -DBUILD_SHARED_LIBS=OFF \
        -DLLAMA_BUILD_TESTS=OFF \
        -DLLAMA_BUILD_EXAMPLES=OFF \
        -DLLAMA_BUILD_SERVER=ON \
        -DLLAMA_CURL=OFF \
        -DLLAMA_BUILD_UI=OFF \
        -DLLAMA_USE_PREBUILT_UI=OFF \
        -DHOST_CXX_COMPILER=/usr/bin/g++ \
        -DGGML_NATIVE=OFF \
        -DGGML_AVX=ON -DGGML_AVX2=ON -DGGML_FMA=ON -DGGML_F16C=ON \
        -DGGML_OPENMP=OFF \
        -DGGML_VULKAN=ON \
        -DVulkan_INCLUDE_DIR=/vk/include \
        -DVulkan_LIBRARY=/vk/lib/libvulkan-1.dll.a \
        -DVulkan_GLSLC_EXECUTABLE=/usr/bin/glslc \
        -DSPIRV-Headers_DIR=/vk/share/cmake/SPIRV-Headers \
        -DCMAKE_EXE_LINKER_FLAGS="-static" \
    && cmake --build /bld --target llama-server -j"$(nproc)"

FROM debian:trixie-slim AS fast-build
RUN apt-get update && apt-get install -y --no-install-recommends \
        git ca-certificates cmake make g++ \
        gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64 \
    && rm -rf /var/lib/apt/lists/*
ARG BERGAMOT_COMMIT=9271618ebbdc5d21ac4dc4df9e72beb7ce644774
ARG MARIAN_COMMIT=2781d735d4a10dca876d61be587afdab2726293c
RUN git clone https://github.com/browsermt/bergamot-translator /bt \
    && cd /bt && git checkout --quiet ${BERGAMOT_COMMIT} \
    && git submodule update --init --recursive \
    && test "$(git rev-parse HEAD)" = "${BERGAMOT_COMMIT}" \
    && test "$(git -C 3rd_party/marian-dev rev-parse HEAD)" = "${MARIAN_COMMIT}"
RUN sed -i 's/^Types: deb$/Types: deb deb-src/' /etc/apt/sources.list.d/debian.sources \
    && apt-get update && apt-get install -y --no-install-recommends dpkg-dev xz-utils \
    && mkdir /pcre && cd /pcre && apt-get source pcre2 \
    && rm -rf /var/lib/apt/lists/* \
    && src=$(find /pcre -maxdepth 1 -mindepth 1 -type d -name 'pcre2-*' | head -1) \
    && test -f "$src/configure" \
    && cp -r "$src" /bt/3rd_party/ssplit-cpp/src/3rd-party/pcre2-10.39
COPY build/mingw-toolchain.cmake /toolchain.cmake
COPY build/fast/fast-translator.cpp /bt/app/fast-translator.cpp
RUN sed -i 's/-Wall; -Werror; -Wextra;/-Wall; -Wno-error; -Wextra;/' /bt/3rd_party/marian-dev/CMakeLists.txt \
    && sed -i -e '25s/^#ifdef _MSC_VER$/#if defined(_MSC_VER) || defined(__MINGW32__)/' -e '83s/^#ifdef _MSC_VER$/#if defined(_MSC_VER) || defined(__MINGW32__)/' -e '3a #ifdef __MINGW32__\n#include <malloc.h>\n#endif' /bt/src/translator/aligned.h \
    && test "$(grep -c __MINGW32__ /bt/src/translator/aligned.h)" = 3 \
    && mkdir -p /shim && for h in Windows Shlobj ShlObj Psapi; do echo "#include <$(echo $h | tr A-Z a-z).h>" > /shim/$h.h; done \
    && printf 'add_executable(fast-translator fast-translator.cpp)\ntarget_link_libraries(fast-translator PRIVATE bergamot-translator)\n' >> /bt/app/CMakeLists.txt \
    && cmake -S /bt -B /bld \
        -DCMAKE_TOOLCHAIN_FILE=/toolchain.cmake \
        -DCMAKE_BUILD_TYPE=Release \
        -DUSE_WASM_COMPATIBLE_SOURCE=ON \
        -DCOMPILE_TESTS=OFF \
        -DGIT_SUBMODULE=OFF \
        -DSSPLIT_USE_INTERNAL_PCRE2=ON \
        -DCMAKE_CXX_STANDARD_INCLUDE_DIRECTORIES=/shim \
        -DCMAKE_C_STANDARD_INCLUDE_DIRECTORIES=/shim \
        -DCMAKE_CXX_FLAGS=-DPCRE2_STATIC \
        -DCMAKE_C_FLAGS=-DPCRE2_STATIC \
        -DCMAKE_CXX_STANDARD_LIBRARIES=-lshlwapi \
        -DCMAKE_EXE_LINKER_FLAGS=-static \
    && cmake --build /bld --target fast-translator -j"$(( $(nproc) < 6 ? $(nproc) : 6 ))" \
    && x86_64-w64-mingw32-strip /bld/app/fast-translator.exe

FROM golang:1.23-bookworm AS client-build
RUN apt-get update && apt-get install -y --no-install-recommends \
        gcc-mingw-w64-x86-64 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY client/go.mod client/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY client/ ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go run ./tools/genicon -o app.ico \
    && x86_64-w64-mingw32-windres -O coff -o resources.syso app.rc
ARG APP_VERSION=
ENV GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -mod=readonly -trimpath \
        -ldflags "-H windowsgui -s -w -extldflags=-static ${APP_VERSION:+-X holdtotype/internal/appid.Version=$APP_VERSION}" \
        -o /out/holdtotype.exe .

FROM debian:bookworm-slim AS sherpa-fetch
RUN apt-get update && apt-get install -y --no-install-recommends \
        curl ca-certificates bzip2 \
    && rm -rf /var/lib/apt/lists/*
ARG SHERPA_ONNX_VERSION=v1.13.6
ARG SHERPA_ONNX_SHA256=388654755fb59c23eefaaab297e2d43660e08c9695cec5b7391ccdd3e274e917
RUN curl -fL --retry 3 -o /s.tar.bz2 \
        "https://github.com/k2-fsa/sherpa-onnx/releases/download/${SHERPA_ONNX_VERSION}/sherpa-onnx-${SHERPA_ONNX_VERSION}-win-x64-static-MT-Release-no-tts.tar.bz2" \
    && echo "${SHERPA_ONNX_SHA256}  /s.tar.bz2" | sha256sum -c - \
    && mkdir /x && tar -xjf /s.tar.bz2 -C /x --strip-components=1 \
    && cp /x/bin/sherpa-onnx-offline-websocket-server.exe /sherpa-server.exe \
    && cp /x/bin/sherpa-onnx-online-websocket-server.exe /sherpa-online-server.exe \
    && rm -rf /x /s.tar.bz2

FROM debian:bookworm-slim AS assets
RUN apt-get update && apt-get install -y --no-install-recommends \
        curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*
ARG WHISPER_MODEL
RUN curl -fL --retry 3 -o /ggml.bin \
        "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-${WHISPER_MODEL}.bin"
COPY client/config.example.json /config.tmpl
RUN sed "s/__MODEL__/ggml-${WHISPER_MODEL}.bin/" /config.tmpl > /config.json

FROM debian:bookworm-slim AS payload
RUN apt-get update && apt-get install -y --no-install-recommends \
        zip \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /p
COPY --from=client-build /out/holdtotype.exe ./
COPY --from=whisper-build /bld/bin/whisper-server.exe ./
COPY --from=llama-build /bld/bin/llama-server.exe ./
COPY --from=llama-gpu-build /bld/bin/llama-server.exe ./llama-server-gpu.exe
COPY --from=fast-build /bld/app/fast-translator.exe ./
COPY --from=sherpa-fetch /sherpa-server.exe /sherpa-online-server.exe ./
COPY --from=assets /config.json ./config.default.json
COPY README.md ./
RUN zip -9 -X payload.zip holdtotype.exe whisper-server.exe llama-server.exe llama-server-gpu.exe fast-translator.exe sherpa-server.exe sherpa-online-server.exe config.default.json README.md

FROM golang:1.23-bookworm AS setup-build
RUN apt-get update && apt-get install -y --no-install-recommends \
        gcc-mingw-w64-x86-64 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY client/go.mod client/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY client/ ./
COPY --from=payload /p/payload.zip setup/payload.zip
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go run ./tools/genicon -o setup/app.ico \
    && cp app.rc app.manifest setup/ \
    && cd setup && x86_64-w64-mingw32-windres -O coff -o resources.syso app.rc
ARG APP_VERSION=
ENV GOOS=windows GOARCH=amd64 CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -mod=readonly -trimpath \
        -ldflags "-H windowsgui -s -w ${APP_VERSION:+-X holdtotype/internal/appid.Version=$APP_VERSION}" \
        -o /out/holdtotype-setup.exe ./setup

FROM golang:1.23-bookworm AS gotest
RUN apt-get update && apt-get install -y --no-install-recommends \
        gcc-mingw-w64-x86-64 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY client/go.mod client/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY client/ ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    touch setup/payload.zip \
    && go test ./internal/... \
    && GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go vet -unsafeptr=false ./...

FROM scratch AS dist-slim
COPY --from=whisper-build /bld/bin/whisper-server.exe /whisper-server.exe
COPY --from=llama-build /bld/bin/llama-server.exe /llama-server.exe
COPY --from=llama-gpu-build /bld/bin/llama-server.exe /llama-server-gpu.exe
COPY --from=fast-build /bld/app/fast-translator.exe /fast-translator.exe
COPY --from=sherpa-fetch /sherpa-server.exe /sherpa-server.exe
COPY --from=sherpa-fetch /sherpa-online-server.exe /sherpa-online-server.exe
COPY --from=client-build /out/holdtotype.exe /holdtotype.exe
COPY --from=setup-build /out/holdtotype-setup.exe /holdtotype-setup.exe
COPY --from=assets /config.json /config.default.json
COPY README.md /README.md

FROM scratch AS dist
ARG WHISPER_MODEL
COPY --from=whisper-build /bld/bin/whisper-server.exe /whisper-server.exe
COPY --from=llama-build /bld/bin/llama-server.exe /llama-server.exe
COPY --from=llama-gpu-build /bld/bin/llama-server.exe /llama-server-gpu.exe
COPY --from=fast-build /bld/app/fast-translator.exe /fast-translator.exe
COPY --from=sherpa-fetch /sherpa-server.exe /sherpa-server.exe
COPY --from=sherpa-fetch /sherpa-online-server.exe /sherpa-online-server.exe
COPY --from=client-build /out/holdtotype.exe /holdtotype.exe
COPY --from=assets /ggml.bin /models/ggml-${WHISPER_MODEL}.bin
COPY --from=assets /config.json /config.default.json
COPY --from=setup-build /out/holdtotype-setup.exe /holdtotype-setup.exe
COPY README.md /README.md
