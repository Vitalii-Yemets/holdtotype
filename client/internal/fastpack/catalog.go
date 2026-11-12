package fastpack

var catalog = []Pack{
	{From: "en", To: "de", Version: "2.1", Files: []File{
		{Kind: "model", Name: "model.ende.intgemm.alphas.bin", Size: 31561787, SHA256: "8df29d9494d19f47fd5d97c6a73474c6f657e9f81c1a607c431d02befdf3810f", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/23db71e7-b6d9-45eb-a47d-0290d7d8ef63.bin"},
		{Kind: "vocab", Name: "vocab.ende.spm", Size: 810073, SHA256: "69f730becafa48e3bb2c244eab66456877c08959a02f2bd5519b5a3088b62f9c", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/261225ea-5a52-455b-981c-7d09c6e6da3c.spm"},
		{Kind: "lex", Name: "lex.50.50.ende.s2t.bin", Size: 4347672, SHA256: "7ed39f1cffbd68a27ddf05bbfe068de2060f1d7e69f1a20e27ae923551dd7393", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/bc072b1a-7749-43f7-9fe0-34a6dff10c4a.bin"},
	}},
	{From: "en", To: "es", Version: "2.1", Files: []File{
		{Kind: "model", Name: "model.enes.intgemm.alphas.bin", Size: 31561787, SHA256: "3b1c399511c01c84c36fae5c0524df44096288efdc8236e182b5c97d7ad2244c", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/a4ba0e94-16de-4058-9a44-5bbbbb3c8640.bin"},
		{Kind: "vocab", Name: "vocab.enes.spm", Size: 816054, SHA256: "5ae254fa9b15aa182e70fd2a6186b1333c63a29a48043a9224c6aa4fcac058ad", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/170634fd-511a-4a28-b723-0a1025c67feb.spm"},
		{Kind: "lex", Name: "lex.50.50.enes.s2t.bin", Size: 4198436, SHA256: "7d51237c0a07027dcd61643cfbbb0f8c48597d19907ef53d2cae9d6bec2cf25c", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/1834a61e-0331-4c4a-bbc0-dda02afa8188.bin"},
	}},
	{From: "en", To: "fr", Version: "2.0", Files: []File{
		{Kind: "model", Name: "model.enfr.intgemm.alphas.bin", Size: 31561787, SHA256: "6322e296d4fecfe395a8d5723da4ec37ecbe6d7613bb1dfcf4b28e2a47498b68", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/e0d7f3ee-163a-4b0a-8e44-8b208b39477d.bin"},
		{Kind: "vocab", Name: "vocab.enfr.spm", Size: 814404, SHA256: "783abf3abe075afdf8d85d233994bef2c3a064e935ab1bed946820aff6ac002a", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/c7c1bcf3-71f6-425a-a6ed-403bd9a0a759.spm"},
		{Kind: "lex", Name: "lex.50.50.enfr.s2t.bin", Size: 4372936, SHA256: "2585ed98d3af0bc949865aedeb390493d591f56870814376e73e4144c41ed059", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/293d1823-ff05-48c9-837b-cff8fd33f197.bin"},
	}},
	{From: "en", To: "it", Version: "2.1", Files: []File{
		{Kind: "model", Name: "model.enit.intgemm.alphas.bin", Size: 31561787, SHA256: "248f47568788ecc351da7e5e07064d4153b4f71e011364ae2c931ffeec4d1cc2", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/29ae0d70-5e37-49e8-8da6-3e2c58da61d1.bin"},
		{Kind: "vocab", Name: "vocab.enit.spm", Size: 812724, SHA256: "3ef0211d4ae6db21440892f180f2019fe2bfc110a330ffa9d2eca9665e4f2bc5", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/d9d5ffa9-b919-4491-9ab6-3dec34459768.spm"},
		{Kind: "lex", Name: "lex.50.50.enit.s2t.bin", Size: 4133192, SHA256: "8b21914804625b2777dae7fdb636eb78f02a3eb8b7bceaa50f29ac740961da93", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/3730ff60-2291-4c87-888d-2dc4593b20db.bin"},
	}},
	{From: "en", To: "pl", Version: "2.1", Files: []File{
		{Kind: "model", Name: "model.enpl.intgemm.alphas.bin", Size: 31561787, SHA256: "b0b59716bcc8f46b8dc6ab5c32a2732acaad42a5a6d9ae9b472526a1bb0d71d4", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/640bd7ef-a977-417f-903c-00b5217cb7de.bin"},
		{Kind: "vocab", Name: "vocab.enpl.spm", Size: 821612, SHA256: "2000f182efa8ce5f69bc3f9e574666b5bc2c0af684eedd32f69d15aa9021ea85", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/dd8a3019-bd2f-405e-b95d-0cce7a29be2d.spm"},
		{Kind: "lex", Name: "lex.50.50.enpl.s2t.bin", Size: 3540052, SHA256: "d7d7b364673f633763c9f8f2e8ee9c7a184e92f4086e8b84abe870a2a8dbe5d5", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/12451e17-4901-4327-b78d-903e9ff7529d.bin"},
	}},
	{From: "en", To: "uk", Version: "2.2", Files: []File{
		{Kind: "model", Name: "model.enuk.intgemm.alphas.bin", Size: 31561787, SHA256: "ed8bb4a811e1c45a259e91052080d25e5c8b4e2df0e59cb1521c04a57855ef97", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/6c69578a-1f77-44ce-993d-4d37c81e081e.bin"},
		{Kind: "vocab", Name: "vocab.enuk.spm", Size: 922077, SHA256: "4dd6a91feff71e745395dcc0d52f885848ca0d843d5f70628ca51b3be1b70b45", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/ffea5811-fcfe-486e-8812-6284bd33ee3a.spm"},
		{Kind: "lex", Name: "lex.50.50.enuk.s2t.bin", Size: 3054364, SHA256: "dbfd4f8dc95da9c575164f7fc95ef00a2fb5f7d9f51b656559f99eeebdb57c87", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/e7b2fd05-abf2-4e30-a612-59bce8a6702c.bin"},
	}},
	{From: "en", To: "ru", Version: "2.1", Files: []File{
		{Kind: "model", Name: "model.enru.intgemm.alphas.bin", Size: 31561787, SHA256: "184cb5cda528eeefc0f75f5d0035d787b71d74af135e3c5608d01ae02ecfb920", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/db6b48bb-b563-4ea1-bbf4-b033fa9cc482.bin"},
		{Kind: "vocab", Name: "vocab.enru.spm", Size: 904455, SHA256: "56ee63e14e8cb926c394242adc3ed7cc602644c3d33058cff2ce2959d52a6258", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/568807dd-9620-4c4d-a0b5-d651d64c0ab0.spm"},
		{Kind: "lex", Name: "lex.50.50.enru.s2t.bin", Size: 2774540, SHA256: "4d91839726b960e70b6d05c53d0cffd16262832b1c0e1ea99d66f412dcc6a239", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/e88cc45c-0e0b-45d3-a957-8d1ba3a9b5aa.bin"},
	}},
	{From: "de", To: "en", Version: "2.0", Files: []File{
		{Kind: "model", Name: "model.deen.intgemm.alphas.bin", Size: 31561787, SHA256: "3e6f7c2c2425d10824797270b382bee718ff34af2cab9308841c82ca46dc6f20", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/f44b1b1b-9df6-4ece-971e-0e5ce96fae54.bin"},
		{Kind: "vocab", Name: "vocab.deen.spm", Size: 810073, SHA256: "69f730becafa48e3bb2c244eab66456877c08959a02f2bd5519b5a3088b62f9c", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/8ad4d93e-21e6-4862-81d5-c1c3a7d0767b.spm"},
		{Kind: "lex", Name: "lex.50.50.deen.s2t.bin", Size: 4945796, SHA256: "113b98460468360cca68c042e1cddf49c4e1931cbb975ed04349c9a3bd607010", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/d0e4efcb-6145-43db-a69e-568904cc2925.bin"},
	}},
	{From: "es", To: "en", Version: "2.0", Files: []File{
		{Kind: "model", Name: "model.esen.intgemm.alphas.bin", Size: 31561787, SHA256: "4aed7734152ae0045d1a69ae49c86cfda18f53c61f90e95e1d1de1c7c7c3b033", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/d7256ea9-7731-4f15-866c-62998eef93b9.bin"},
		{Kind: "vocab", Name: "vocab.esen.spm", Size: 816054, SHA256: "5ae254fa9b15aa182e70fd2a6186b1333c63a29a48043a9224c6aa4fcac058ad", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/8695cc3d-7897-4c29-808b-bc983a93d1cd.spm"},
		{Kind: "lex", Name: "lex.50.50.esen.s2t.bin", Size: 4636248, SHA256: "e2610211d3b9577d012638fe7e7e74ed7b4b708ce96b9e792e67c282a6492daa", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/a55a9cf2-8a8c-4942-a528-91e93d081746.bin"},
	}},
	{From: "fr", To: "en", Version: "2.0", Files: []File{
		{Kind: "model", Name: "model.fren.intgemm.alphas.bin", Size: 31561787, SHA256: "15f997bc0d13808b0b0fbd0786e684a3c8a52adcd8071844b76123fdacbf2b90", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/8ef78b1c-df6c-49fc-8cae-cf8c0c2a6630.bin"},
		{Kind: "vocab", Name: "vocab.fren.spm", Size: 814404, SHA256: "783abf3abe075afdf8d85d233994bef2c3a064e935ab1bed946820aff6ac002a", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/11ecb9da-60dd-4b39-aa12-27883099294b.spm"},
		{Kind: "lex", Name: "lex.50.50.fren.s2t.bin", Size: 4824120, SHA256: "87c6752ea908f5f0347c10ac0cf7d80d9c2f4f20c81c90168f3e8230b56d4440", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/c4d99496-81bf-480f-b0bb-3aebacdcbe8e.bin"},
	}},
	{From: "it", To: "en", Version: "2.0", Files: []File{
		{Kind: "model", Name: "model.iten.intgemm.alphas.bin", Size: 31561787, SHA256: "21b70978ce2f3b4da7a06b5de86a09abe3acd30b9eee1b2ebb3582b9bad790bf", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/a6e3de84-6612-4546-a4ef-b658ee23ca08.bin"},
		{Kind: "vocab", Name: "vocab.iten.spm", Size: 812724, SHA256: "3ef0211d4ae6db21440892f180f2019fe2bfc110a330ffa9d2eca9665e4f2bc5", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/51b5a437-9daa-481e-a9bc-311f2b3f490c.spm"},
		{Kind: "lex", Name: "lex.50.50.iten.s2t.bin", Size: 4713616, SHA256: "0700d6d70b30490e68ff9deca0b45b80310745b491855bde3a13e1692d0cbce1", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/122ac310-0f7b-40d5-b49d-bcee81077c85.bin"},
	}},
	{From: "pl", To: "en", Version: "2.0", Files: []File{
		{Kind: "model", Name: "model.plen.intgemm.alphas.bin", Size: 31561787, SHA256: "06dbab3baa3bbceab746bc6a3bba807878f38536382c03acddf59d3edf3b5d94", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/dd8a7ab8-b12f-4ad4-8827-0922510df232.bin"},
		{Kind: "vocab", Name: "vocab.plen.spm", Size: 821612, SHA256: "2000f182efa8ce5f69bc3f9e574666b5bc2c0af684eedd32f69d15aa9021ea85", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/47b0f4cd-7d3c-4cf5-b7b4-0aecfbdc51f6.spm"},
		{Kind: "lex", Name: "lex.50.50.plen.s2t.bin", Size: 4964208, SHA256: "b295dfefa7022a62e105276e8c39fdfd67bffe588f4fd78583edd4e04ee0cb3b", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/9a072be8-d900-4062-a2ab-c53c8eb8e619.bin"},
	}},
	{From: "uk", To: "en", Version: "1.1", Files: []File{
		{Kind: "model", Name: "model.uken.intgemm.alphas.bin", Size: 17141051, SHA256: "2e7b8bc9d2569d8b4ae244e964cdc343a0aa3735cc83ecb664be84a37cb45f48", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/31e03e41-af47-4de2-9420-f6fd206e0114.bin"},
		{Kind: "vocab", Name: "vocab.uken.spm", Size: 886410, SHA256: "8b51480dfc4ca2d7a08d2bfdaa9e57353ca97e6cf3d5f9aba4c0179d3665bfc0", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/3155c3f0-b486-40f3-bfdd-1b3b734cc171.spm"},
		{Kind: "lex", Name: "lex.50.50.uken.s2t.bin", Size: 3914120, SHA256: "32b5e0df493671a223cb5cb32430181cbb7fe6b5562c91dbc3d35410d37e5ff6", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/3e01a41b-9384-4bfc-98b9-031df4f9ee5e.bin"},
	}},
	{From: "ru", To: "en", Version: "2.1", Files: []File{
		{Kind: "model", Name: "model.ruen.intgemm.alphas.bin", Size: 31561787, SHA256: "aa0efa194b3dbd2da129ba0b8db8ef139126cb0914e57072f62061294746f1e2", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/20260901210056--9781ee2a-d15e-4fe8-893d-687193373fff--37b65106-c6cd-4ed2-a28e-3acd72725429.bin"},
		{Kind: "vocab", Name: "vocab.ruen.spm", Size: 952371, SHA256: "de8af5f87e231ee4016cec4bb482e6b584cb1ca34dcb624fc7b8398f25bb35df", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/20260901210058--222363ca-4a98-48a9-9162-65acdffc2572--ba370dc8-709d-42cd-b68c-83b912bd3e38.spm"},
		{Kind: "lex", Name: "lex.50.50.ruen.s2t.bin", Size: 4397908, SHA256: "c944e753f298bbdf6506388de7257fb48cbd12c83e2807e2c284fc323771895b", URL: "https://firefox-settings-attachments.cdn.mozilla.net/main-workspace/translations-models/20260901210100--fd808a7b-0081-427d-83a5-7043bdd97a87--215ca362-0436-495d-b89d-542187b696cb.bin"},
	}},
}
