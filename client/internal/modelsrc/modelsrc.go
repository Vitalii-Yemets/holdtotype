package modelsrc

const (
	hub     = "https://huggingface.co/"
	resolve = "/resolve/"

	Whisper  = hub + "ggerganov/whisper.cpp" + resolve + "5359861c739e955e79d9a303bcbc70fb988958b1/"
	GigaAMv3 = hub + "csukuangfj/sherpa-onnx-nemo-transducer-punct-giga-am-v3-russian-2025-12-16" + resolve + "a6039be7cee829a9044a69ac0ebaf1c191217c97/"
	GigaAMv2 = hub + "csukuangfj/sherpa-onnx-nemo-transducer-giga-am-v2-russian-2025-04-19" + resolve + "d7002de64f758983ac1146695e5de5a381c36785/"
	Parakeet = hub + "csukuangfj/sherpa-onnx-nemo-parakeet-tdt-0.6b-v3-int8" + resolve + "2bda32ec70b097a55adaa07d9a7173915b43cc78/"
	Nemotron = hub + "Masterx/sherpa-onnx-nemotron-3.5-asr-streaming-0.6b-560ms-2026-06-11" + resolve + "3f1f6020667b151d0b6820cc4b8ae4f42b6d35fc/"
	Canary   = hub + "csukuangfj/sherpa-onnx-nemo-canary-180m-flash-en-es-de-fr-int8" + resolve + "9077164e0d3dd1d5353743e89ceaa1d3a770838c/"
	Qwen3ASR = hub + "csukuangfj2/sherpa-onnx-qwen3-asr-0.6B-int8-2026-03-25" + resolve + "68818b2313fe77bd06f6a7c5068ff3ef59d02b8a/"
)

func All() map[string]string {
	return map[string]string{
		"Whisper": Whisper, "GigaAMv3": GigaAMv3, "GigaAMv2": GigaAMv2, "Parakeet": Parakeet,
		"Nemotron": Nemotron, "Canary": Canary, "Qwen3ASR": Qwen3ASR,
	}
}
