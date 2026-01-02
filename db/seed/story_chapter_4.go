package seed

import (
	"encoding/json"
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/internal/domain/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedChapter4(db *gorm.DB, vocabMap map[string]uuid.UUID) error {
	slog.Info("seeding chapter 4...")

	var chapter entity.Chapter
	err := db.Where("order_index = ?", 4).First(&chapter).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			chapter = entity.Chapter{
				Title:         "Ngadepi Juragan Cengkeh",
				Description:   "Ujian pungkasan. Andi bakal entuk restu apa malah kena penthung tongkate Pak Broto?",
				CoverImageURL: "chapters/ch4_cover.webp",
				OrderIndex:    4,
			}
			if err := db.Create(&chapter).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	andi := func(exp string) charData { return charData{Name: "Andi", Img: "chars/andi_" + exp + ".webp"} }
	sekar := func(exp string) charData { return charData{Name: "Sekar", Img: "chars/sekar_" + exp + ".webp"} }
	pakbroto := func(exp string) charData { return charData{Name: "Pak Broto", Img: "chars/pakbroto_" + exp + ".webp"} }

	slidesData := []slideData{
		// intro
		{Key: "1", Speaker: "Narator", BgImg: "bg/halaman_pak_broto.webp", Content: "Langit ing Tulungagung katingal mendhung. Griya Joglo ageng ing {ngarsanipun} Andi krasa kadosdene kraton.", NextSlideKey: "2", VocabKeys: []string{"ngarsanipun"}},
		{Key: "2", Speaker: "Narator", BgImg: "bg/halaman_pak_broto.webp", Content: "{Wancinipun} mbuktekaken asil pasinaon. Andi ngatur ambegan, nyiapaken mental kangge ngadhepi pacoban pungkasan.", NextSlideKey: "3", VocabKeys: []string{"wancinipun"}},
		{Key: "3", Speaker: "Andi", BgImg: "bg/halaman_pak_broto.webp", Characters: []charData{andi("nervous")}, Content: "(Ngadek ngarep lawang jati sing gedhe lan ukir-ukiran) (_Bismillah... Eling pesene Pakdhe Joyo. Aja grusa-grusu. Aja ndredeg. Duh, tapi sikile lemes._)", NextSlideKey: "4"},
		{Key: "4", Speaker: "Sekar", BgImg: "bg/halaman_pak_broto.webp", Characters: []charData{andi("nervous"), sekar("worried")}, Content: "Mas Andi! Sampun {dugi}? Monggo, Bapak sampun {ngrantos} ing lebet.", NextSlideKey: "5", VocabKeys: []string{"dugi", "ngrantos"}},

		// choice 1
		{
			Key: "5", Speaker: "Andi", BgImg: "bg/halaman_pak_broto.webp", Characters: []charData{andi("nervous"), sekar("worried")},
			Content: "Inggih Dik. (Andi mlaku nyedaki lawang utama sing menga sithik)",
			Choices: []choiceSeedData{
				{Text: "Permisi...", NextSlideKey: "6a", MoodImpact: 0},
				{Text: "{Kula nuwun}...", NextSlideKey: "6b", MoodImpact: 1},
				{Text: "Assalamualaikum Pak Broto!", NextSlideKey: "6c", MoodImpact: -2},
			},
			VocabKeys: []string{"kula nuwun"},
		},

		// branch choice 1
		{Key: "6a", Speaker: "Pak Broto", BgImg: "bg/halaman_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "(Saka njero, swarane abot lan adhem) Sapa iku? Kok ora duwe tata krama! Durung dipersilahkan kok wis 'permisi'.", NextSlideKey: "7a"},
		{Key: "7a", Speaker: "Andi", BgImg: "bg/halaman_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "(_Waduh, salah start iki._)", NextSlideKey: "8"},

		{Key: "6b", Speaker: "Pak Broto", BgImg: "bg/halaman_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "(Saka njero, wibawa) Monggo... {Mlebet}.", NextSlideKey: "7b", VocabKeys: []string{"mlebet"}},
		{Key: "7b", Speaker: "Andi", BgImg: "bg/halaman_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("neutral")}, Content: "(_Alhamdulillah, lawang dibukak._)", NextSlideKey: "8"},

		{Key: "6c", Speaker: "Pak Broto", BgImg: "bg/halaman_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("angry")}, Content: "(Mendengus banter) Bengok-bengok kaya neng alas. Iki omah, udu lapangan!", NextSlideKey: "7c"},
		{Key: "7c", Speaker: "Andi", BgImg: "bg/halaman_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("angry")}, Content: "(_Mati aku, Pak Broto langsung bad mood._)", NextSlideKey: "8"},

		// merge path
		{Key: "8", Speaker: "Narator", BgImg: "bg/ruang_tamu_pak_broto.webp", Content: "Andi {mlebet} ing ruang tamu. Prabot jati kuno lan lukisan jaran nambah kesan wibawa ing ruangan menika.", NextSlideKey: "9", VocabKeys: []string{"mlebet"}},
		{Key: "9", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("intimidating")}, Content: "(Lungguh ing kursi jati, nyekel tongkat komando, natah Andi saka ndhuwur nganti ngisor tanpa kedhep) ...", NextSlideKey: "10"},
		{Key: "10", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("intimidating")}, Content: "(_Waduh, matane kaya elang arep nyaut pitik. Aku pitike._)", NextSlideKey: "11"},
		{Key: "11", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("intimidating")}, Content: "(Mlaku nyedaki Pak Broto, mbungkuk sithik) {Sugeng} {sonten}, Pak.", NextSlideKey: "12", VocabKeys: []string{"sugeng", "sonten"}},
		{Key: "12", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "{Lenggah}.", NextSlideKey: "13", VocabKeys: []string{"lenggah"}},

		// choice 2
		{
			Key: "13", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")},
			Content: "(Kudu lungguh ing kursi kayu sing atos).",
			Choices: []choiceSeedData{
				{Text: "Lungguh senderan ben rileks.", NextSlideKey: "14a", MoodImpact: -1},
				{Text: "Lungguh tegap, tangan Ngapurancang.", NextSlideKey: "14b", MoodImpact: 1},
				{Text: "Lungguh mbungkuk banget.", NextSlideKey: "14c", MoodImpact: 0},
			},
		},

		// branch choice 2
		{Key: "14a", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "Enak men lungguhmu. Kaya juragan wae. (Nyindir pedhes)", NextSlideKey: "15"},
		{Key: "14b", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "(Mung nglirik tangan Andi, manthuk sithik) Hmm.", NextSlideKey: "15"},
		{Key: "14c", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "Aja mbungkuk-mbungkuk nemen. Kowe arep lamaran apa arep ngemis?", NextSlideKey: "15"},

		// merge path
		{Key: "15", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Dadi... iki sing jenenge Andi?", NextSlideKey: "16"},
		{Key: "16", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Inggih Pak, {leres}. {Kula} Andi.", NextSlideKey: "17", VocabKeys: []string{"kula", "leres"}},
		{Key: "17", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Bakul mie? (Nadane datar, ora ngenyek tapi ngetes mental).", NextSlideKey: "18"},

		// choice 3
		{
			Key: "18", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")},
			Content: "(_Waduh, 'Bakul Mie' kok krasa nylekit ya._)",
			Choices: []choiceSeedData{
				{Text: "Inggih Pak, namung bakul mie.", NextSlideKey: "19a", MoodImpact: 0},
				{Text: "Inggih Pak, {kula} {sadeyan} mie.", NextSlideKey: "19b", MoodImpact: 1},
				{Text: "CEO Warmindo Pak.", NextSlideKey: "19c", MoodImpact: -2},
			},
			VocabKeys: []string{"kula", "sadeyan"},
		},

		// branch choice 3
		{Key: "19a", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "Yen kowe dhewe ora bangga karo gaweanmu, piye aku iso percaya kowe iso nguripi anakku?", NextSlideKey: "20"},
		{Key: "19b", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("neutral")}, Content: "{Sadeyan} mie... Hmm. Usaha halal. Ora masalah.", NextSlideKey: "20", VocabKeys: []string{"sadeyan"}},
		{Key: "19c", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("angry")}, Content: "CEO? Gayane selangit. Isih warung cilik wae wis umuk.", NextSlideKey: "20"},

		{Key: "20", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Sekar iku anakku wedok siji-sijine. Tak gedhekne karo disiplin lan kecukupan. Kowe wani njamin apa gawe masa depane?", NextSlideKey: "21"},

		// choice 4
		{
			Key: "21", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("intimidating")},
			Content: "(Pak Broto condong menyang ngarep, natah tajem).",
			Choices: []choiceSeedData{
				{Text: "Wah, asil {kula} atusan yuta Pak!", NextSlideKey: "22a", MoodImpact: -1},
				{Text: "{Kula} janji {badhe} ngebahagiakne Sekar.", NextSlideKey: "22b", MoodImpact: 0},
				{Text: "Insyaallah {cekap} Pak. {Kula} {badhe} ikhtiar.", NextSlideKey: "22c", MoodImpact: 1},
			},
			VocabKeys: []string{"kula", "cekap", "badhe"},
		},

		// branch choice 4
		{Key: "22a", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "Hah, isih enom kok wis pamer harta. Duwit iso entek sak kedhepan mata.", NextSlideKey: "23"},
		{Key: "22b", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Bahagia iku butuh bukti, udu janji manis.", NextSlideKey: "23"},
		{Key: "22c", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("neutral")}, Content: "Hmm... Jawabanmu apik. Wong lanang pancen kudu wani tanggung jawab lan usaha.", NextSlideKey: "23"},

		// merge path
		{Key: "23", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("neutral")}, Content: "Oh inggih, Bapak. Menika {wonten} {sakedhik} tandha tresna {saking} Surabaya. (Nyerahke bungkusan Batik).", NextSlideKey: "24", VocabKeys: []string{"wonten", "sakedhik", "saking"}},
		{Key: "24", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "(Nampani tanpa ngomong, mbukak bungkus pelan-pelan).", NextSlideKey: "25"},
		{Key: "25", Speaker: "Narator", BgImg: "bg/ruang_tamu_pak_broto.webp", Content: "Swasana hening malih. Namung swara kertas krekek-krekek ingkang kepireng. Andi nahan ambegan.", NextSlideKey: "26"},
		{Key: "26", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "(Ndelok motif batike) Wahyu Tumurun... Sogan.", NextSlideKey: "27"},
		{Key: "27", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "Sapa sing milihne iki? Kowe?", NextSlideKey: "28"},

		// choice 5
		{
			Key: "28", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")},
			Content: "(Kudu njawab jujur nanging sopan).",
			Choices: []choiceSeedData{
				{Text: "Inggih Pak, {kula} mireng Bapak {remen}.", NextSlideKey: "29a", MoodImpact: 1},
				{Text: "Batik Sogan, Pak. Jarene apik gawe sampeyan.", NextSlideKey: "29b", MoodImpact: -2},
				{Text: "Niki Batik larang lho Pak, Sutra asli.", NextSlideKey: "29c", MoodImpact: -1},
			},
			VocabKeys: []string{"kula", "remen"},
		},

		// branch choice 5
		{Key: "29a", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("happy")}, Content: "(Mesem tipis banget, meh ora ketara) Bagus. Kowe ngerti seleraku. Jarang anak muda saiki ngerti batik.", NextSlideKey: "28_2"},
		{Key: "29b", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("angry")}, Content: "(Mendelik) '{Sampeyan}'? Kowe nganggep aku kancamu ta? Batik apik dadi ora aji merga basamu kasar.", NextSlideKey: "28_2", VocabKeys: []string{"sampeyan"}},
		{Key: "29c", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "Aku ora takon regane. Barang larang nek sing ngekek ora ikhlas ya ora ana gunane.", NextSlideKey: "28_2"},

		{Key: "28_2", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Sekar! Gawekno kopi gawe Andi. Kopi ireng, aja gulo.", NextSlideKey: "29_2"},
		{Key: "29_2", Speaker: "Sekar", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), sekar("happy"), pakbroto("neutral")}, Content: "(Saka mburi lawang, mesem seneng) Inggih Pak.", NextSlideKey: "30"},

		// merge path
		{Key: "30", Speaker: "Narator", BgImg: "bg/ruang_tamu_pak_broto.webp", Content: "Sekar medal mbeta kopi. Ambunipun sedhep, nanging uabipun taksih kemebul panas.", NextSlideKey: "31"},
		{Key: "31", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Diunjuk kopine.", NextSlideKey: "32"},
		{Key: "32", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Inggih Pak. (Andi nyekel gelas, panas!)", NextSlideKey: "33"},

		// choice 6
		{
			Key: "33", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")},
			Content: "(_Waduh, uap e isih kemebul. Iki nek tak ombe lambeku melepuh. Tapi Pak Broto wis ngakon._)",
			Choices: []choiceSeedData{
				{Text: "Langsung sruput.", NextSlideKey: "34a", MoodImpact: -1},
				{Text: "Sekedap Pak, {ngrantos} {asrep}.", NextSlideKey: "34b", MoodImpact: 0},
				{Text: "Nyebul kopi pelan, lagi diombe sithik.", NextSlideKey: "34c", MoodImpact: 1},
			},
			VocabKeys: []string{"ngrantos", "asrep"},
		},

		// branch choice 6
		{Key: "34a", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "(Nyruput) Adhuh! (Kepanasen, kopi kutah sithik).", NextSlideKey: "35a"},
		{Key: "35a", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "(Geleng-geleng) Grusa-grusu. Wong lanang iku kudu tenang. Kopi panas kok diombe.", NextSlideKey: "36"},
		{Key: "34b", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "Kok mung didelok? Ora doyan kopi ta?", NextSlideKey: "36"},
		{Key: "34c", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "(Nyebul alon, nyruput sithik) Sruput...", NextSlideKey: "35c"},
		{Key: "35c", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("happy")}, Content: "Ya ngono. Alon-alon asal kelakon. Urip bebojoan ya ngono, kudu sabar.", NextSlideKey: "36"},

		// merge path
		{Key: "36", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Ndi. Jujur wae. Aku wis krungu akeh babagan kowe saka Sekar.", NextSlideKey: "37"},
		{Key: "37", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Awale aku ragu. Bocah Suroboyoan, bakul mie, urakan.", NextSlideKey: "38"},
		{Key: "38", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "(Nunduk) Ngapunten Pak...", NextSlideKey: "39"},
		{Key: "39", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("neutral")}, Content: "Nanging dina iki, kowe teka mrene. Kowe berusaha nggawe Basa Krama senajan ilatmu kaku.", NextSlideKey: "40"},
		{Key: "40", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("happy")}, Content: "Kowe wani ngadepi aku dhewean. Iku jenenge Lakon Lanang. Saiki pertanyaan terakhir. Yen aku ngijini kowe rabi karo Sekar, apa janjimu?", NextSlideKey: "41"},

		// choice 7
		{
			Key: "41", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")},
			Content: "(Andi sumpah)",
			Choices: []choiceSeedData{
				{Text: "{Kula} janji Sekar {mboten} bakal keliren.", NextSlideKey: "42a", MoodImpact: 0},
				{Text: "{Kula} janji {badhe} njagi lan nuntun Sekar.", NextSlideKey: "42b", MoodImpact: 2},
				{Text: "Aku janji gak bakal nglarani atine.", NextSlideKey: "42c", MoodImpact: -5},
			},
			VocabKeys: []string{"kula", "mboten", "badhe"},
		},

		// branch choice 7
		{Key: "42a", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Materi iku nomer sekian. Sing penting iku ketentreman.", NextSlideKey: "43"},
		{Key: "42b", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("happy")}, Content: "Iku jawaban sing tak tunggu. Imam kudu iso nuntun makmum.", NextSlideKey: "43"},
		{Key: "42c", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("angry")}, Content: "Hadeh... wis arep bener malah kepleset ngoko neng pungkasan. Sinau maneh!", NextSlideKey: "43"},

		// merge path
		{Key: "43", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("neutral"), pakbroto("neutral")}, Content: "Yowes. Aku titip anakku. Aja disia-sia ne. Kapan wong tuwamu iso mrene?", NextSlideKey: "44"},
		{Key: "44", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("happy")}, Content: "(Kaget seneng, ngangkat sirah) {Estu} Pak? Bapak {paring} restu? Inggih Pak! Minggu ngajeng {kula} ajak Bapak Ibu mriki!", NextSlideKey: "45", VocabKeys: []string{"estu", "paring", "kula"}},
		{Key: "45", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("neutral")}, Content: "Iyo. Tapi eling, nek kowe wani macem-macem, tongkatku iki sing ngomong.", NextSlideKey: "46"},
		{Key: "46", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("neutral")}, Content: "{Mboten} Pak! {Kula} janji.", NextSlideKey: "47", VocabKeys: []string{"mboten", "kula"}},
		{Key: "47", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("neutral")}, Content: "Wis kana, balika. Wis sore.", NextSlideKey: "48"},

		// choice 8
		{
			Key: "48", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("neutral")},
			Content: "(Ngadeg, raine sumringah) (_Alhamdulillah! Sukses rek!_)",
			Choices: []choiceSeedData{
				{Text: "Suwun Pak, aku balik sek.", NextSlideKey: "49a", MoodImpact: -1},
				{Text: "Matur nuwun Pak, {kula} {nyuwun} {pamit}.", NextSlideKey: "49b", MoodImpact: 1},
				{Text: "Nggih Pak, dadah.", NextSlideKey: "49c", MoodImpact: -1},
			},
			VocabKeys: []string{"kula", "nyuwun"},
		},

		// branch choice 8
		{Key: "49a", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("angry")}, Content: "(Gedheg-gedheg) Isih kudu akeh sinau bocah iki...", NextSlideKey: "50"},
		{Key: "49b", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy"), pakbroto("happy")}, Content: "Ati-ati. Salam gawe wong tuwamu.", NextSlideKey: "50"},
		{Key: "49c", Speaker: "Pak Broto", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("nervous"), pakbroto("angry")}, Content: "Dadah? Kowe pikir aku kancamu?", NextSlideKey: "50"},

		// closing
		{Key: "50", Speaker: "Narator", BgImg: "bg/ruang_tamu_pak_broto.webp", Content: "Pungkasane, Andi kasil ngalahake rasa {ajrih}ipun lan pikantuk restu saking Pak Broto.", NextSlideKey: "51", VocabKeys: []string{"ajrih"}},
		{Key: "51", Speaker: "Narator", BgImg: "bg/ruang_tamu_pak_broto.webp", Content: "Dedemen amargi sampurna, nanging amargi purun mbudidaya lan ngurmati tiyang sanes.", NextSlideKey: "52"},
		{Key: "52", Speaker: "Andi", BgImg: "bg/ruang_tamu_pak_broto.webp", Characters: []charData{andi("happy")}, Content: "(_Maturnuwun Gusti... Akhire rabi!_)", NextSlideKey: "53"},
		{Key: "53", Speaker: "Narator", BgImg: "bg/wedding_venue.webp", Content: "Lakon Sowan sampun purna. Andi lan Sekar miwiti lembaran enggal kanthi restu lan kabagyan.", NextSlideKey: "54"},
		{Key: "54", Speaker: "Narator", BgImg: "bg/wedding_venue.webp", Content: "TAMAT."},
	}

	realIDs := make(map[string]uuid.UUID)

	for _, d := range slidesData {
		slide := entity.Slide{
			ChapterID:          chapter.ID,
			SpeakerName:        d.Speaker,
			Content:            d.Content,
			BackgroundImageURL: d.BgImg,
		}

		var charJson []map[string]interface{}
		for _, c := range d.Characters {
			isActive := (c.Name == d.Speaker)
			charJson = append(charJson, map[string]interface{}{
				"name":      c.Name,
				"image_url": c.Img,
				"is_active": isActive,
			})
		}
		if len(charJson) > 0 {
			bytes, _ := json.Marshal(charJson)
			slide.Characters = types.JSONB(bytes)
		} else {
			slide.Characters = types.JSONB("[]")
		}

		for _, vKey := range d.VocabKeys {
			if vID, ok := vocabMap[vKey]; ok {
				slide.Vocabularies = append(slide.Vocabularies, entity.Dictionary{ID: vID})
			}
		}

		if err := db.Create(&slide).Error; err != nil {
			return err
		}

		realIDs[d.Key] = slide.ID
	}

	for _, d := range slidesData {
		updates := map[string]interface{}{}

		if d.NextSlideKey != "" {
			if nextRealID, ok := realIDs[d.NextSlideKey]; ok {
				updates["next_slide_id"] = nextRealID
			}
		}

		if len(d.Choices) > 0 {
			updates["choices"] = makeChoicesWithRealIDs(d.Choices, realIDs)
		}

		if len(updates) > 0 {
			if err := db.Model(&entity.Slide{}).Where("id = ?", realIDs[d.Key]).Updates(updates).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
