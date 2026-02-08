-- Импорт старых групп из JSON (DedStack, ТИМА ДОДГЕРА РОДГЕРА, Beyond, и т.д.)
-- Запуск: sqlite3 /path/to/dayzsmartcf.db < seed_old_groups.sql
-- Или из папки backend: sqlite3 ../../dayzsmartcf.db < scripts/seed_old_groups.sql

-- 1) Игроки: добавляем по cftools_id, если ещё нет (display_name из никнейма/алиаса)
INSERT OR IGNORE INTO players (cftools_id, display_name) VALUES
('5f9dac373642e730d425a767', ':flag_jp:Urajiku:flag_jp:'),
('5d4284b5a1db3a062cf0b149', 'VENOM'),
('63a86b88f2ce6c273766e958', 'so bad in this place'),
('5e8ce20f7fa53497b65f6437', '[SBU] Izmailov'),
('635719d89e4b77ec93f15710', '[:ninja: :flag_ru:/PEEK] BJJ'),
('60e83956de09345bb1f40fcc', 'Rak prostati'''),
('615623c14013fdc984a23458', 'GGorynich'),
('63a5bbf89e4b77ec9365d569', 'Kystik'),
('5e0e56dc2b764260c5afaa0c', 'midnight'),
('5ff0a348cd785fbcfb715f9c', 'LB.PeepXD'),
('6280f71000a0aa43c8b8e4dc', 'Dvasoska'),
('5ff5852ff380d5d07b99af2f', 'ovoshnoy'),
('5ff5852ff380d5d07b99af32', 'MsrH'),
('63a7fff4c869e9a047844209', 'RED'),
('637d4e4290d48f5870f81294', 'evp. botik<3'),
('61c36d8d3de0a4ce0aaf24e0', 'rururu'),
('63c1b160d7410c3ad5d3e832', 'unc:man_factory_worker:  Rodjomba'),
('626a56a59ed3f821cfb83614', '[UA/PEEK] DAHbIK'),
('5bcb835cc221d513de8f1cb3', 'ShaMan'),
('5f1751ffb31e0ab6456165b7', '*ssk/ Polka'),
('602d6046b5abb416d4ab4258', 'Berserk:flag_ru:'),
('5c8fe0ada1db3a45c246a6bb', 'MAFIOZZzi1'),
('5e651001289728a4affa27c2', 'Tengen?'),
('610326847c78ada1393c1634', 'odinrazpokaazivau'),
('5c6c344fa1db3a2e9a9eba38', 'Like Nastya'),
('61061c266d6cf2108d58911c', 'Kayota (2)'),
('64550f38089abe445a48bf5a', 'KaDuHa'),
('5cdd705fa1db3a6c607a0cef', 'CEBEP'),
('61fae89d305432c43acfb508', 'DarkSide :smiling_imp: '),
('64136f241311f898ee095022', 'SSS+'),
('609e672b9ad1701bd0969695', '.fairytale*'),
('61f54ba3cd978bff0dbe613d', 'Beyond'),
('61fe756dcd978bff0dee635d', 'Survivor (3)'),
('607eb0e03b2c5f27d48e0bd2', '23Green23'),
('641b56c41311f898ee394b53', 'SHEFF-KZ'),
('5eca5982ec2ead5662bd4e6a', 'BadassKOTG'),
('5c2d07e8a1db3a19c1984e50', 'Fluttershy'),
('61c58d0773d0b5d02bd12ff3', 'MyRaVeY:ant:'),
('6204f425e9ab54c993d78add', 'enemy.Mura'),
('5e8cb7bb93953e1a0455138b', 'Snake'),
('62414fceb02293033817a186', 'enemy.Yandex_Luftwaffe'),
('5beac603c221d5392d756a1e', 'enemy.MORTA'),
('62865c23f689f5ed49d7c6e7', 'enemy.mister'),
('63a81e8789d42ed78961b622', 'CTblD'),
('645f6f8a2d7cf421d8920a8a', 'ELLIE'),
('63ac4fc1f806a39f159d01ee', 'enemy.yushka666'),
('5fd13e259954c7e6fa71d106', 'boring_Materialist'),
('5fe5fb999706616f18293d0b', 'Vampir'),
('5ef3c56be4c470e1e27a9146', 'TTuBHou_Pa3HoC'),
('5ff447ac46152a080d927b36', 'EHoT (2)'),
('63a7305a9e4b77ec936cfc2e', 'enemy.Exylie'),
('627bac59d523d0d6369e004e', '[KM] // B+B'),
('64135be96f6ad524660b8963', 'Survivor'),
('6450d86ed61a1ac40d5d06ef', 'KM//Nekit Roiz:black_heart:'),
('610fb3f71f2bde02c756daea', '[MON] Drozd (2)'),
('6485f8608e6cf5721b4ad61a', 'Reanim :ambulance:'),
('5e4d7026248374e1679ff77c', '[MRKV] Iskander'),
('64a2cce05530a1d82a2d887d', 'Gatsby'),
('5c7ac993a1db3a0e5c2ec394', 'evp. Crymikal:flag_cz: (2)'),
('649c4c72b8bb176569c33097', '[MRKV] Marmyshka'),
('63a87a1cd24d061b55b9eb18', 'Zeexdi (2)'),
('62a8cf6f98599c3b0b1dab38', 'Demon'),
('63bdf8889142c3610b978a5a', 'luvv//akia'),
('5f8ea453cd2c3ecf3f55f351', 'Survivor (2)'),
('645f7a86bdde226343741231', 'chz./Sokolnikov'),
('6233762b0bc7061a63f4e648', 'ebz.rockabye:peach:'),
('5c2df4dea1db3a19c1997f3e', 'evp. :snake:'),
('615540de2500244784c250fd', 'evp. Exact'),
('610451e106cb44634d53de41', 'Opezdol (2)'),
('64ac5b239e47055167409191', '3060 // sliyanie'),
('61e095f8ea590deacf34313d', 'bu'),
('6575d90699756a607538b10b', 'Prezident Koshmara'),
('64846d7cb8bb1765691f64e1', '75 // suport'),
('5c380950a1db3a0a315f828a', '75//fy'),
('5e78eb75a7974f93decfb5e7', 'Are$'),
('6242322a013470036742b79d', 'KakZheYaEby'),
('607f28913b2c5f27d49167bf', 'kr0nya'),
('61a2bfd97ac61cf4576c4f13', '! kakoy chudesnyy den'' chtoby podarit'' tsvetok :wilted_'),
('6521812176fd477a5e9fa5ac', 'Maniiak :x_ray:'),
('62b0bb2a6c4d5ab3434c033f', 'Survivor (5)'),
('5d4e83b5a1db3a062c2b8845', '[AUT] Creon'),
('5ceee55ea1db3a6c60c25081', 'ceweqe.'),
('5e43fe06fb2c4374aa538e42', 'AUT Egomaniac'),
('5fc0d085942e405b93a46502', '[CBTM] Kykysya'),
('6605a8fce0689c0731c6f517', '[CBTM]RoDiPiT'),
('63de7396a14038d638f863f6', '=R@Ge='),
('628a73b85bd379ff29dedcec', 'nolly'),
('60d6574eb1166a1de6c5c777', 'insff'),
('5bf42d3ec221d54ead5c2567', '+w absolute???'),
('60f419316535a50395132804', 'merc.Goshan'),
('61362418e841211778cdce5f', 'merc.FI3IK'),
('61d0cd2b1595d6147de011c9', 'merc. :strawberry:'),
('5e2704743205107c045320b7', 'merc.Bear '),
('63dd0101ce3d6bc4eedb8b12', 'merc.koza v tazike'),
('60d32ec97f1b279fa528b0f3', './/b4by?'),
('63b7b43d874140ef5deb878b', 'thn'),
('62790f6400a0aa43c892087c', 'merc.WaiKiKi'),
('63b4c8dfe98b700bfd881f27', 'merc.Navfrik'),
('601da8ca5eb9f32e6379079e', 'merc.Fantomas'),
('60b00d61ff4ab494e855bc43', 'Bogatyri | LiGhTninG'),
('61066bbcfc2700266560ad76', 'Kametra ??'),
('5c9dd290a1db3a29d885649c', ':circus_tent: Karma'),
('620c86e41c0b219b59d0c9d2', 'Bogatyri | Yari4ek'),
('5eaaab70ccc616b8dd779603', '19ruslan85-rus'),
('635e51484eb59a23643fc5fb', 'Rumka_Vodki'),
('632628ac88911a97faa75dd3', '[dal_bano4ky]Dora ne Dura'),
('647a3edbc398ed78b5371fe4', 'UwagA// Vladislavch'),
('607701273b2c5f27d45f5e3e', 'UwagA//GENAKRAHABOR -_-?'),
('6387c7ef1eadb3ebb8dc039d', '[Pussy] 0.5 kd razboynik'),
('6360d4527a3e561d3c4f3c97', '[dal_bano4ky]GEORGIANDATRA'),
('63aff555fa11ccb22f7c6115', 'UwagA//-ArSiK-'),
('5f11b7181808829dad58c62a', 'UwagA// MeowMeeN'),
('61d594d5320aeb3c10811787', 'Mini Pekka'),
('5ff09c4989f4847720a05997', '[BR] Bydka'),
('61c7906000ba65e2f0ce712d', 'Serega kemper'),
('64ca516b8eabcc78746c2b74', '[MS}Baldwin'),
('5de66c12b2dea1bd18d6cced', 'Olegjan'),
('65db9c974995cf27b311c2b4', 'Tonnabatona'),
('627645e05bd379ff29843548', 'ls//:foggy:'),
('63825c9ac8420dfd65d82cde', 'Koja'),
('62deb0c43724eab5c4e63a51', 'GENA NA GENA DERJI'),
('63b9bdd789d42ed789a92ff5', '[BBBBBRRRRRRRRR] :flag_de:'),
('6264326b2389f7feffbd729f', 'wm.MERSI'),
('63fd1c7e6f6ad524668c17be', 'qul1nek1'),
('6317705e3a686698347a4099', '[-RS-] Hitler :smiling_imp:'),
('61b38a8c01181d3b3974254b', 'night.kennedy:raccoon:'),
('5ecab5a85932c93f07ce6dc5', '-1'),
('5fd68b796dd25b8ed523d3a7', 'frozy:pouring_liquid:'),
('5cdbf9dea1db3a6c6075bcf8', 'UwagA// Y-3'),
('5f9ad0917fa875dd4f01bf24', 'Vasia Brigadir'),
('5fe634c85291f84b5c2aad7b', 'Medoed'),
('61dc29f620266e6589ccf82e', '8{'),
('625ff00c6795981466b07c9c', 'DVOEIIINIK'),
('5cffba9ca1db3a6c600bbbde', 'Flyaga'),
('6037d19949fddcd7f52a3a4b', ''),
('605cbf2930d00ebea3eb1b97', ''),
('5e963e991b72135e89711bb3', 'ruja mavpa'),
('62928f021b8bea0cf4cb961b', 'Aza:purple_heart:'),
('605218bef0f81123e5e2c664', 'AndroMeda'),
('5d8e47c8d2c98ecb78dcf405', 'NASVAY .// Sperma V smetane (2)'),
('63573ae4c8420dfd65181183', 'TEHb'),
('6301cea40572c55439032a94', 'Niller Baks'),
('5ef89f5ccf81d911f885cb8a', 'Volshebnik Krisolov'),
('627004a767ed6fe625e52f7a', 'Tequila :flag_by: (2)'),
('61cac38dea590deacfd9fbb0', '//nrxt pisya v gorle'),
('5f0dc1b51808829dad47ee23', 'Zonin?'),
('65d8b878c111b10cb7e8348d', '[ALFA]whiskey'),
('5e403ffd41cc3ad1014ba640', 'global01');

-- 2) Группы (id из старого формата, чтобы не пересекаться с будущими авто-id)
INSERT OR REPLACE INTO groups (id, name, created_at, updated_at) VALUES
(1784102140, 'DedStack', datetime('now'), datetime('now')),
(4798335658, 'ТИМА ДОДГЕРА РОДГЕРА', datetime('now'), datetime('now')),
(4159800265, 'Beyond', datetime('now'), datetime('now')),
(9789187781, '[EN-Y]', datetime('now'), datetime('now')),
(4404718656, 'simurg', datetime('now'), datetime('now')),
(8228909082, 'petrovka 5x', datetime('now'), datetime('now')),
(7516361386, 'p2pishka', datetime('now'), datetime('now')),
(4946983739, 'ВЫШЕ ДОЛЬНИКА ШТУКА', datetime('now'), datetime('now')),
(56233942, 'UWAGA', datetime('now'), datetime('now')),
(2957250865, 'merc', datetime('now'), datetime('now')),
(1462141834, 'cheaters', datetime('now'), datetime('now')),
(3350914050, 'UwagA', datetime('now'), datetime('now')),
(1372237691, 'Aza ruja mavpa', datetime('now'), datetime('now')),
(6536016762, 'Тэкила', datetime('now'), datetime('now'));

-- 3) Участники групп (group_id, player_id из players по cftools_id, alias)
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'Urajiku' FROM players WHERE cftools_id = '5f9dac373642e730d425a767';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'venom' FROM players WHERE cftools_id = '5d4284b5a1db3a062cf0b149';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'weakness' FROM players WHERE cftools_id = '63a86b88f2ce6c273766e958';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'puzo' FROM players WHERE cftools_id = '5e8ce20f7fa53497b65f6437';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'BJJ' FROM players WHERE cftools_id = '635719d89e4b77ec93f15710';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'spyyakinesis' FROM players WHERE cftools_id = '60e83956de09345bb1f40fcc';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'GGorinych' FROM players WHERE cftools_id = '615623c14013fdc984a23458';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'Enforcer' FROM players WHERE cftools_id = '63a5bbf89e4b77ec9365d569';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'midnight' FROM players WHERE cftools_id = '5e0e56dc2b764260c5afaa0c';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'LB.PeepXD' FROM players WHERE cftools_id = '5ff0a348cd785fbcfb715f9c';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'DVASOSKA' FROM players WHERE cftools_id = '6280f71000a0aa43c8b8e4dc';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'mrsh' FROM players WHERE cftools_id = '5ff5852ff380d5d07b99af2f';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1784102140, id, 'ovoshnoy' FROM players WHERE cftools_id = '5ff5852ff380d5d07b99af32';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'ReadHead' FROM players WHERE cftools_id = '63a7fff4c869e9a047844209';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'Borik19' FROM players WHERE cftools_id = '637d4e4290d48f5870f81294';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'Evs' FROM players WHERE cftools_id = '61c36d8d3de0a4ce0aaf24e0';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'dodgerRodger' FROM players WHERE cftools_id = '63c1b160d7410c3ad5d3e832';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'Данык' FROM players WHERE cftools_id = '626a56a59ed3f821cfb83614';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'ShaMan' FROM players WHERE cftools_id = '5bcb835cc221d513de8f1cb3';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'Dem' FROM players WHERE cftools_id = '5f1751ffb31e0ab6456165b7';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'Berserk' FROM players WHERE cftools_id = '602d6046b5abb416d4ab4258';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'Maffiozzi1' FROM players WHERE cftools_id = '5c8fe0ada1db3a45c246a6bb';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'Tengen' FROM players WHERE cftools_id = '5e651001289728a4affa27c2';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'random' FROM players WHERE cftools_id = '610326847c78ada1393c1634';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'Valkin Ded' FROM players WHERE cftools_id = '5c6c344fa1db3a2e9a9eba38';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'donk666' FROM players WHERE cftools_id = '61061c266d6cf2108d58911c';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'Kaduha_BLR' FROM players WHERE cftools_id = '64550f38089abe445a48bf5a';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'JlaJla' FROM players WHERE cftools_id = '5cdd705fa1db3a6c607a0cef';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, '76561198163001878' FROM players WHERE cftools_id = '61fae89d305432c43acfb508';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'jkzop7' FROM players WHERE cftools_id = '64136f241311f898ee095022';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4798335658, id, 'fairytail' FROM players WHERE cftools_id = '609e672b9ad1701bd0969695';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4159800265, id, 'Beyond' FROM players WHERE cftools_id = '61f54ba3cd978bff0dbe613d';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4159800265, id, 'qqs' FROM players WHERE cftools_id = '61fe756dcd978bff0dee635d';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4159800265, id, '23green23' FROM players WHERE cftools_id = '607eb0e03b2c5f27d48e0bd2';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4159800265, id, 'SHEFF' FROM players WHERE cftools_id = '641b56c41311f898ee394b53';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4159800265, id, 'BadassKOTG' FROM players WHERE cftools_id = '5eca5982ec2ead5662bd4e6a';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4159800265, id, 'Moon Baboon' FROM players WHERE cftools_id = '5c2d07e8a1db3a19c1984e50';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4159800265, id, 'mypavey' FROM players WHERE cftools_id = '61c58d0773d0b5d02bd12ff3';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'mura' FROM players WHERE cftools_id = '6204f425e9ab54c993d78add';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'Snake' FROM players WHERE cftools_id = '5e8cb7bb93953e1a0455138b';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'Gnidling' FROM players WHERE cftools_id = '62414fceb02293033817a186';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'morta' FROM players WHERE cftools_id = '5beac603c221d5392d756a1e';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'mister' FROM players WHERE cftools_id = '62865c23f689f5ed49d7c6e7';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'pisya_v_pope' FROM players WHERE cftools_id = '63a81e8789d42ed78961b622';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'Ellie' FROM players WHERE cftools_id = '645f6f8a2d7cf421d8920a8a';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'finish' FROM players WHERE cftools_id = '63ac4fc1f806a39f159d01ee';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'boring_Materialist' FROM players WHERE cftools_id = '5fd13e259954c7e6fa71d106';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'Vampir' FROM players WHERE cftools_id = '5fe5fb999706616f18293d0b';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'Ani4ka_2007' FROM players WHERE cftools_id = '5ef3c56be4c470e1e27a9146';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'EHOT' FROM players WHERE cftools_id = '5ff447ac46152a080d927b36';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 9789187781, id, 'Exylie' FROM players WHERE cftools_id = '63a7305a9e4b77ec936cfc2e';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, '[OGOROD]_ Bonch+Bruevich' FROM players WHERE cftools_id = '627bac59d523d0d6369e004e';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, '[WbZ] Vostok...' FROM players WHERE cftools_id = '64135be96f6ad524660b8963';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'Nekit Roiz' FROM players WHERE cftools_id = '6450d86ed61a1ac40d5d06ef';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'Drozd' FROM players WHERE cftools_id = '610fb3f71f2bde02c756daea';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'Reanim' FROM players WHERE cftools_id = '6485f8608e6cf5721b4ad61a';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, '[MRKV] Iskander' FROM players WHERE cftools_id = '5e4d7026248374e1679ff77c';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'gatsby' FROM players WHERE cftools_id = '64a2cce05530a1d82a2d887d';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'EvP // Crymikal' FROM players WHERE cftools_id = '5c7ac993a1db3a0e5c2ec394';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, '[MRKV] peredoz' FROM players WHERE cftools_id = '649c4c72b8bb176569c33097';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, ' [MRKV] Zeexdi' FROM players WHERE cftools_id = '63a87a1cd24d061b55b9eb18';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'Efiop' FROM players WHERE cftools_id = '62a8cf6f98599c3b0b1dab38';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'psychoakia' FROM players WHERE cftools_id = '63bdf8889142c3610b978a5a';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'EVP' FROM players WHERE cftools_id = '5f8ea453cd2c3ecf3f55f351';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'EVP' FROM players WHERE cftools_id = '645f7a86bdde226343741231';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'WbZ' FROM players WHERE cftools_id = '6233762b0bc7061a63f4e648';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'Viper from ua' FROM players WHERE cftools_id = '5c2df4dea1db3a19c1997f3e';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'MrMikroname' FROM players WHERE cftools_id = '615540de2500244784c250fd';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4404718656, id, 'Opezdol' FROM players WHERE cftools_id = '610451e106cb44634d53de41';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 8228909082, id, '3060' FROM players WHERE cftools_id = '64ac5b239e47055167409191';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 8228909082, id, 'bu' FROM players WHERE cftools_id = '61e095f8ea590deacf34313d';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 8228909082, id, 'prezident' FROM players WHERE cftools_id = '6575d90699756a607538b10b';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 8228909082, id, '75 // suport' FROM players WHERE cftools_id = '64846d7cb8bb1765691f64e1';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 8228909082, id, 'fy' FROM players WHERE cftools_id = '5c380950a1db3a0a315f828a';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 7516361386, id, 'Ares' FROM players WHERE cftools_id = '5e78eb75a7974f93decfb5e7';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 7516361386, id, 'Tamelark' FROM players WHERE cftools_id = '6242322a013470036742b79d';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 7516361386, id, 'Kr0nya' FROM players WHERE cftools_id = '607f28913b2c5f27d49167bf';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 7516361386, id, '! kakoy chudesnyy den'' chtoby podarit'' tsvetok :wilted_' FROM players WHERE cftools_id = '61a2bfd97ac61cf4576c4f13';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 7516361386, id, 'KakZheYaSlab' FROM players WHERE cftools_id = '6521812176fd477a5e9fa5ac';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 7516361386, id, 'p2pishka' FROM players WHERE cftools_id = '62b0bb2a6c4d5ab3434c033f';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4946983739, id, 'q' FROM players WHERE cftools_id = '5d4e83b5a1db3a062c2b8845';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4946983739, id, 'q' FROM players WHERE cftools_id = '5ceee55ea1db3a6c60c25081';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 4946983739, id, 'q' FROM players WHERE cftools_id = '5e43fe06fb2c4374aa538e42';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 56233942, id, 'kukysya' FROM players WHERE cftools_id = '5fc0d085942e405b93a46502';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 56233942, id, 'rodipit' FROM players WHERE cftools_id = '6605a8fce0689c0731c6f517';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 56233942, id, '=R@Ge=' FROM players WHERE cftools_id = '63de7396a14038d638f863f6';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 56233942, id, 'nolly' FROM players WHERE cftools_id = '628a73b85bd379ff29dedcec';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 56233942, id, 'insff' FROM players WHERE cftools_id = '60d6574eb1166a1de6c5c777';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 56233942, id, '+w absolute' FROM players WHERE cftools_id = '5bf42d3ec221d54ead5c2567';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'merc.Goshan' FROM players WHERE cftools_id = '60f419316535a50395132804';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'merc.FI3IK' FROM players WHERE cftools_id = '61362418e841211778cdce5f';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'merc. :strawberry:' FROM players WHERE cftools_id = '61d0cd2b1595d6147de011c9';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'merc.Bear' FROM players WHERE cftools_id = '5e2704743205107c045320b7';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'koza v tazike' FROM players WHERE cftools_id = '63dd0101ce3d6bc4eedb8b12';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'warlock' FROM players WHERE cftools_id = '60d32ec97f1b279fa528b0f3';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, ' merc.WaiKiKi' FROM players WHERE cftools_id = '62790f6400a0aa43c892087c';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'merc.Navfrik' FROM players WHERE cftools_id = '63b4c8dfe98b700bfd881f27';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'merc.Fantomas' FROM players WHERE cftools_id = '601da8ca5eb9f32e6379079e';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'Bogatyri | LiGhTninG' FROM players WHERE cftools_id = '60b00d61ff4ab494e855bc43';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'Kametra' FROM players WHERE cftools_id = '61066bbcfc2700266560ad76';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'Harassment' FROM players WHERE cftools_id = '5c9dd290a1db3a29d885649c';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 2957250865, id, 'Bogatyri | Yari4ek' FROM players WHERE cftools_id = '620c86e41c0b219b59d0c9d2';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1462141834, id, 'mamylek' FROM players WHERE cftools_id = '5eaaab70ccc616b8dd779603';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1462141834, id, 'rumka_vodki' FROM players WHERE cftools_id = '635e51484eb59a23643fc5fb';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'dora' FROM players WHERE cftools_id = '632628ac88911a97faa75dd3';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'vladislavch' FROM players WHERE cftools_id = '647a3edbc398ed78b5371fe4';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'gena' FROM players WHERE cftools_id = '607701273b2c5f27d45f5e3e';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, '0.5kd' FROM players WHERE cftools_id = '6387c7ef1eadb3ebb8dc039d';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'andatra' FROM players WHERE cftools_id = '6360d4527a3e561d3c4f3c97';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, '-ArSiK-' FROM players WHERE cftools_id = '63aff555fa11ccb22f7c6115';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'meowmeen' FROM players WHERE cftools_id = '5f11b7181808829dad58c62a';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'Mini Pekka' FROM players WHERE cftools_id = '61d594d5320aeb3c10811787';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'Elektri4eskaya Bydka' FROM players WHERE cftools_id = '5ff09c4989f4847720a05997';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, ' [MS] Flowby :crown:' FROM players WHERE cftools_id = '61a2bfd97ac61cf4576c4f13';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, '{MS}Serg2.0' FROM players WHERE cftools_id = '61c7906000ba65e2f0ce712d';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, '[MS]Baldwin' FROM players WHERE cftools_id = '64ca516b8eabcc78746c2b74';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'Olegjan' FROM players WHERE cftools_id = '5de66c12b2dea1bd18d6cced';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'eby_kaba4ki' FROM players WHERE cftools_id = '65db9c974995cf27b311c2b4';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'ls//:foggy:' FROM players WHERE cftools_id = '627645e05bd379ff29843548';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'Koja' FROM players WHERE cftools_id = '63825c9ac8420dfd65d82cde';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'GENA NA GENA DERJI' FROM players WHERE cftools_id = '62deb0c43724eab5c4e63a51';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, '[BBBBBRRRRRRRRR] :flag_de:' FROM players WHERE cftools_id = '63b9bdd789d42ed789a92ff5';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'wm.MERSI' FROM players WHERE cftools_id = '6264326b2389f7feffbd729f';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, '76561199114783066' FROM players WHERE cftools_id = '63fd1c7e6f6ad524668c17be';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'AMORI FATI' FROM players WHERE cftools_id = '6317705e3a686698347a4099';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, '1' FROM players WHERE cftools_id = '61b38a8c01181d3b3974254b';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, '1' FROM players WHERE cftools_id = '5ecab5a85932c93f07ce6dc5';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'Frozy' FROM players WHERE cftools_id = '5fd68b796dd25b8ed523d3a7';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'Y-3' FROM players WHERE cftools_id = '5cdbf9dea1db3a6c6075bcf8';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'vasia' FROM players WHERE cftools_id = '5f9ad0917fa875dd4f01bf24';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'olejkaa' FROM players WHERE cftools_id = '5fe634c85291f84b5c2aad7b';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'serg' FROM players WHERE cftools_id = '61dc29f620266e6589ccf82e';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'DVOEIIINIK' FROM players WHERE cftools_id = '625ff00c6795981466b07c9c';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'Flyaga' FROM players WHERE cftools_id = '5cffba9ca1db3a6c600bbbde';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'Jrat_hosh' FROM players WHERE cftools_id = '6037d19949fddcd7f52a3a4b';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 3350914050, id, 'Ya_ne_goloden' FROM players WHERE cftools_id = '605cbf2930d00ebea3eb1b97';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1372237691, id, 'ruja mavpa' FROM players WHERE cftools_id = '5e963e991b72135e89711bb3';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1372237691, id, 'Aza:purple_heart:' FROM players WHERE cftools_id = '62928f021b8bea0cf4cb961b';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 1372237691, id, 'andromeda' FROM players WHERE cftools_id = '605218bef0f81123e5e2c664';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 6536016762, id, 'каша' FROM players WHERE cftools_id = '5d8e47c8d2c98ecb78dcf405';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 6536016762, id, 'хируко' FROM players WHERE cftools_id = '63573ae4c8420dfd65181183';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 6536016762, id, 'эзбукс' FROM players WHERE cftools_id = '6301cea40572c55439032a94';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 6536016762, id, 'замонолит' FROM players WHERE cftools_id = '5ef89f5ccf81d911f885cb8a';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 6536016762, id, 'тэкила' FROM players WHERE cftools_id = '627004a767ed6fe625e52f7a';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 6536016762, id, 'lincor' FROM players WHERE cftools_id = '61cac38dea590deacfd9fbb0';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 6536016762, id, 'пара' FROM players WHERE cftools_id = '5f0dc1b51808829dad47ee23';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 6536016762, id, 'бобр' FROM players WHERE cftools_id = '65d8b878c111b10cb7e8348d';
INSERT OR IGNORE INTO group_members (group_id, player_id, alias) SELECT 6536016762, id, 'global01' FROM players WHERE cftools_id = '5e403ffd41cc3ad1014ba640';
