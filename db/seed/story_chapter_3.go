package seed

import (
	"encoding/json"
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/internal/domain/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedChapter3(db *gorm.DB, vocabMap map[string]uuid.UUID) error {
	slog.Info("seeding chapter 3...")

	var chapter entity.Chapter
	err := db.Where("order_index = ?", 3).First(&chapter).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			chapter = entity.Chapter{
				Title:         "Golek Gawan Sowan",
				Description:   "Andi kudu pinter milih batik lan nawar rego ngadepi Bu Tejo sing galak.",
				CoverImageURL: "chapters/ch3_cover.webp",
				OrderIndex:    3,
			}
			if err := db.Create(&chapter).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	andi := func(exp string) charData { return charData{Name: "Andi", Img: "chars/andi_" + exp + ".webp"} }
	butejo := func(exp string) charData { return charData{Name: "Bu Tejo", Img: "chars/butejo_" + exp + ".webp"} }

	slidesData := []slideData{
		// intro
		{Key: "1", Speaker: "Narator", BgImg: "bg/toko_batik.webp", Content: "Pakdhe Joyo ngutus Andi tumuju dhateng Pasar Besar. Ananging, papan ingkang dipuntuju sanes toko sembarangan.", NextSlideKey: "2", VocabKeys: []string{"sowan"}},
		{Key: "2", Speaker: "Narator", BgImg: "bg/toko_batik.webp", Content: "Toko {menika} namanipun 'Batik Lestari', {kagunganipun} Bu Tejo. Piyantun Solo ingkang sampun dangu {wonten} Surabaya, nanging kenceng anggenipun ngugemi tata krama.", NextSlideKey: "3", VocabKeys: []string{"menika", "kagungan", "wonten"}},
		{Key: "3", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous")}, Content: "(Mandheg ngarep toko sing akeh hiasan wayang lan kain jarik) (_Waduh, ambune dupa menyan. Iki toko batik apa dukun? Pakdhe Joyo pancen aneh-aneh wae._)", NextSlideKey: "4"},
		{Key: "4", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral")}, Content: "(Mlebu toko, Bu Tejo lagi sibuk ngetung duit neng kalkulator, ora noleh) ...", NextSlideKey: "5"},

		// choice 1
		{
			Key: "5", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("neutral")},
			Content: "(Dicuekin, toleh-toleh bingung) (_Waduh, dicuekin rek. Kudu piye iki? Bengok apa nunggu?_)",
			Choices: []choiceSeedData{
				{Text: "Bu! {Tumbas} Bu! Halo!", NextSlideKey: "6a", MoodImpact: -1},
				{Text: "Ngenteni kanthi sabar.", NextSlideKey: "6b", MoodImpact: 1},
				{Text: "Ehem! Nuwun sewu Bu...", NextSlideKey: "6c", MoodImpact: 0},
			},
			VocabKeys: []string{"tumbas"},
		},

		// branch choice 1
		{Key: "6a", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("shocked")}, Content: "(Kaget, kalkulatore gigal) Astaghfirullah! Mas! Jantungku meh coplok! Sabar {sakedhik} napa!", NextSlideKey: "7a", VocabKeys: []string{"sakedhik"}},
		{Key: "7a", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("angry")}, Content: "Eh, ngapunten Bu. Kesusu.", NextSlideKey: "8"},

		{Key: "6b", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "(Meneng ngadeg kaya patung.)", NextSlideKey: "7b"},
		{Key: "7b", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "(Akhire sadar, noleh) Lho, wis ket {kalawau} ta Le? Kok ora muni? Tak pikir manekin anyar.", NextSlideKey: "7b_2", VocabKeys: []string{"kalawau"}},
		{Key: "7b_2", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("neutral")}, Content: "Hehe, {boten} Bu. Kula ngrantos Ibu salse.", NextSlideKey: "8", VocabKeys: []string{"boten"}},

		{Key: "6c", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Sik Mas, nanggung. Kurang {sakedhik} iki itungane.", NextSlideKey: "7c", VocabKeys: []string{"sakedhik"}},
		{Key: "7c", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Inggih Bu.", NextSlideKey: "8"},

		// merge path
		{Key: "8", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("happy")}, Content: "Iyo wis. Arep golek apa? Batik cap apa tulis? Apa nggolek jodoh? (Ngguyu cekikikan)", NextSlideKey: "9"},
		{Key: "9", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Wah, Ibu niki saged mawon. Kula {madosi} batik kangge calon mertua, Bu.", NextSlideKey: "10", VocabKeys: []string{"madosi"}},
		{Key: "10", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("neutral")}, Content: "Woalah, arep '{sowan}' ta? Pantesan raine tegang kaya klambi durung disetrika.", NextSlideKey: "11", VocabKeys: []string{"sowan"}},
		{Key: "11", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Calon mertuane wong {pundi}? Nek wong Suroboyoan, tak jupukne sing motif Suro lan Boyo sing garang.", NextSlideKey: "12", VocabKeys: []string{"pundi"}},
		{Key: "12", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Sanes Bu. Tiyang Tulungagung. Priyayi sepuh.", NextSlideKey: "13"},
		{Key: "13", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("shocked")}, Content: "Waduh! Tulungagung? Priyayi? Iku selerane kudu 'Kelas Berat' Le. Gak iso sembarangan.", NextSlideKey: "14"},
		{Key: "14", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "(Munggah kursi, njupuk kain neng rak paling dhuwur) Iki lho. Batik Sogan motif Wahyu Tumurun. Iki jimat ampuh gawe ngluluhke ati calon mertua.", NextSlideKey: "15"},
		{Key: "15", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("neutral")}, Content: "(Ndelok kain sing ditunjuk) Wah, apik tenan. Warnane kalem. Kula {badhe} ningali ingkang... (Bingung nunjuk kain sing endi)", NextSlideKey: "16", VocabKeys: []string{"badhe"}},

		// choice 2
		{
			Key: "16", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")},
			Content: "(Bingung nunjuk kain sing endi)",
			Choices: []choiceSeedData{
				{Text: "Sing {niku} Bu.", NextSlideKey: "17a", MoodImpact: 0},
				{Text: "Sing kuwi Bu.", NextSlideKey: "17b", MoodImpact: -1},
				{Text: "Ingkang {menika} Bu.", NextSlideKey: "17c", MoodImpact: 1},
			},
			VocabKeys: []string{"niku", "menika"},
		},

		// branch choice 2
		{Key: "17a", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "'{Niku}'? Sing {pundi}? Sing kiwa apa tengen? Nunjuk sing jelas ta Le.", NextSlideKey: "18", VocabKeys: []string{"niku", "pundi"}},
		{Key: "17b", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("angry")}, Content: "Heh? 'Kuwi'? Mas, iki batik tulis alus, nyebut e sing alus {sakedhik}.", NextSlideKey: "18", VocabKeys: []string{"sakedhik"}},
		{Key: "17c", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Oh, Ingkang {menika}? Pilihan cerdas. Iki favorit para pejabat jaman mbiyen.", NextSlideKey: "18", VocabKeys: []string{"menika"}},

		// merge path
		{Key: "18", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Wis, dicepeng dhisik. Alus ta? Sutra asli iki.", NextSlideKey: "19"},
		{Key: "19", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("neutral")}, Content: "Inggih Bu, alus sanget. Kados pipine... eh, {boten}.", NextSlideKey: "20", VocabKeys: []string{"boten"}},
		{Key: "20", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("happy")}, Content: "Hayo, mikir anake sapa? Fokus Mas, fokus! Pripun? Cocok {boten}?", NextSlideKey: "21", VocabKeys: []string{"boten"}},

		// choice 3
		{
			Key: "21", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")},
			Content: "(Andi arep milih)",
			Choices: []choiceSeedData{
				{Text: "{Kula} nyuwun sing niki mawon.", NextSlideKey: "22a", MoodImpact: 1},
				{Text: "Aku njaluk sing iki wae.", NextSlideKey: "22b", MoodImpact: -1},
				{Text: "{Kula} purun sing niki.", NextSlideKey: "22c", MoodImpact: 0},
			},
			VocabKeys: []string{"kula"},
		},

		// branch choice 3
		{Key: "22a", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Sing niki nggih? Leres. {Sampeyan} pinter milih.", NextSlideKey: "23", VocabKeys: []string{"sampeyan"}},
		{Key: "22b", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("angry")}, Content: "(Nyentak) Aku? Njaluk? {Tumbas} Mas, sanes njaluk! Basane dijaga!", NextSlideKey: "23", VocabKeys: []string{"tumbas"}},
		{Key: "22c", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "'Purun'? Inggih monggo. Basa {sampeyan} lucu, kaku kaya Londo sinau Jawa.", NextSlideKey: "23", VocabKeys: []string{"sampeyan"}},

		// merge path
		{Key: "23", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Niki {reginipun} {pinten} Bu?", NextSlideKey: "24", VocabKeys: []string{"regi", "pinten"}},
		{Key: "24", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Murah. Gawe awakmu sing ganteng, 350 {ewu} wae.", NextSlideKey: "25", VocabKeys: []string{"ewu"}},
		{Key: "25", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("neutral")}, Content: "(_Mak glodhak! 350 {ewu}? Iso ora mangan seminggu aku iki._)", NextSlideKey: "26", VocabKeys: []string{"ewu"}},
		{Key: "26", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("neutral")}, Content: "Waduh Bu... Kok {awis} sanget? {Kula} niki nembe merintis usaha. Lho, rega nggawa rupa Mas. Tapi yowes, tawarana. Tak rungokne.", NextSlideKey: "27", VocabKeys: []string{"awis", "kula"}},

		// choice 4
		{
			Key: "27", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")},
			Content: "(Andi nyiapake tawaran)",
			Choices: []choiceSeedData{
				{Text: "200 {ewu} lah Bu! Pas!", NextSlideKey: "28a", MoodImpact: -1},
				{Text: "Napa {mboten} saged kirang, Bu? 250 {ewu} {pripun}?", NextSlideKey: "28b", MoodImpact: 1},
				{Text: "Larang men Bu! Toko sebelah luwih murah!", NextSlideKey: "28c", MoodImpact: -2},
			},
			VocabKeys: []string{"ewu", "mboten", "pripun"},
		},

		// branch choice 4
		{Key: "28a", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("angry")}, Content: "200? Kowe arep ngerampok aku ta? Gak oleh! Modal benange wae wis larang.", NextSlideKey: "29"},
		{Key: "28b", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Hmm... 250 ya? Jan-jane rugi aku. Tapi berhubung kowe sopan lan arep 'berjuang' demi cinta... Yowes lah. Bungkus!", NextSlideKey: "29"},
		{Key: "28c", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("angry")}, Content: "(Nesu, kain ditarik maneh) Yowes tuku o neng toko sebelah kana! Aja balik mrene! (Andi panik njaluk sepura)", NextSlideKey: "29"},

		// merge path
		{Key: "29", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Alhamdulillah. Matur nuwun Bu Tejo.", NextSlideKey: "30"},
		{Key: "30", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("neutral")}, Content: "Sekalian Mas, bungkuse sing apik ya. Supaya Camer {remen}. Tambah jajanan sisan apa ora?", NextSlideKey: "31", VocabKeys: []string{"remen"}},
		{Key: "31", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Oh inggih. Gethuk lindri niki sae.", NextSlideKey: "32"},
		{Key: "32", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "{Badhe} {mundhut} {pinten} kotak?", NextSlideKey: "33", VocabKeys: []string{"badhe", "mundhut", "pinten"}},

		// choice 5
		{
			Key: "33", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")},
			Content: "(Andi arep tuku 2 kotak)",
			Choices: []choiceSeedData{
				{Text: "Loro.", NextSlideKey: "34a", MoodImpact: 0},
				{Text: "{Kalih}.", NextSlideKey: "34b", MoodImpact: 1},
				{Text: "{Kalih} {atus}.", NextSlideKey: "34c", MoodImpact: -1},
			},
			VocabKeys: []string{"kalih", "atus"},
		},

		// branch choice 5
		{Key: "34a", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "'Loro'? . '{Kalih}' ngono lho.", NextSlideKey: "35", VocabKeys: []string{"kalih"}},
		{Key: "34b", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "{Kalih} kotak nggih. Siap.", NextSlideKey: "35", VocabKeys: []string{"kalih"}},
		{Key: "34c", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("shocked"), butejo("shocked")}, Content: "Hah?! {Kalih} {atus} (200) kotak? {Badhe} slametan tiyang sak kampung ta Mas?", NextSlideKey: "35", VocabKeys: []string{"kalih", "atus", "badhe"}},

		// choice 6
		{
			Key: "35", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")},
			Content: "Dadi totale: Batik 250 {ewu} + Gethuk 50 {ewu}. Kabeh dadi...",
			Choices: []choiceSeedData{
				{Text: "Telung {atus} {ewu}.", NextSlideKey: "36a", MoodImpact: 0},
				{Text: "{Tiga} {atus} {ewu}.", NextSlideKey: "36b", MoodImpact: 1},
				{Text: "Telu {atus} {ewu}.", NextSlideKey: "36c", MoodImpact: -1},
			},
			VocabKeys: []string{"ewu", "atus", "tiga"},
		},

		// branch choice 6
		{Key: "36a", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "{Tiga} {atus} Mas. 'Telung' menika ngoko.", NextSlideKey: "37", VocabKeys: []string{"tiga", "atus"}},
		{Key: "36b", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Leres. {Tiga} {atus} {ewu} pas.", NextSlideKey: "37", VocabKeys: []string{"tiga", "atus", "ewu"}},
		{Key: "36c", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("neutral")}, Content: "Telu {atus}? Basa {pundi} menika Mas? Ingkang leres {Tiga} {Atus}.", NextSlideKey: "37", VocabKeys: []string{"pundi", "tiga", "atus"}},

		// merge path
		{Key: "37", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Niki {artanipun}, Bu. Pas nggih.", NextSlideKey: "38", VocabKeys: []string{"arta"}},
		{Key: "38", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Matur nuwun. Eh Mas, kok {sampeyan} getol banget sinau Basa? Padahal arek Suroboyo.", NextSlideKey: "39", VocabKeys: []string{"sampeyan"}},
		{Key: "39", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Inggih Bu. Demi calon mertua, {kula} purun sinau. Piyantunipun Priyayi Sepuh soalipun.", NextSlideKey: "40", VocabKeys: []string{"kula"}},
		{Key: "40", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Apik. Mas, tak kandhani ya. Aku iki wis dodolan batik puluhan taun. Wis apal watake wong Jawa, utamane Priyayi.", NextSlideKey: "41"},
		{Key: "41", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "{Sampeyan} nek ngadepi Priyayi sepuh, kuncine mung {setunggal}. Napa niku Bu? {Kula} butuh contekan.", NextSlideKey: "42", VocabKeys: []string{"sampeyan", "setunggal", "kula"}},

		// choice 7
		{
			Key: "42", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")},
			Content: "{Kula} butuh contekan.",
			Choices: []choiceSeedData{
				{Text: "Bapak senengane opo?", NextSlideKey: "43a", MoodImpact: 0},
				{Text: "{Kersanipun} Bapak {menika} kados pundi?", NextSlideKey: "43b", MoodImpact: 1},
				{Text: "Bapak {remenipun} napa?", NextSlideKey: "43c", MoodImpact: 0},
			},
			VocabKeys: []string{"kula", "kersa", "menika", "remen"},
		},

		// branch choice 7
		{Key: "43a", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "Duh, basamu lho Mas. Isih pating pecotot. '{Kersanipun}' ngono lho.", NextSlideKey: "44", VocabKeys: []string{"kersa"}},
		{Key: "43b", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Nah, pinter. Priyayi iku seneng wong sing 'Genah'. Nek ditakoni, jawabe sing mantep. Aja plin-plan.", NextSlideKey: "44"},
		{Key: "43c", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")}, Content: "'{Remenipun}' wis bener, tapi kurang alus {sakedhik} nek kangge Priyayi sepuh.", NextSlideKey: "44", VocabKeys: []string{"remen", "sakedhik"}},

		// merge path
		{Key: "44", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("teaching")}, Content: "Lan siji maneh... Aja sok sugih. Priyayi jaman semono luwih ngregani 'Unggah-ungguh' timbang unggah-unggahan bondo.", NextSlideKey: "45"},
		{Key: "45", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "(Manthuk-manthuk) (_Ngena banget pesene Ibu iki._) Inggih Bu, matur nuwun.", NextSlideKey: "46"},
		{Key: "46", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Wis, tak dongakne lancar ya Le. Mugi-mugi calon mertuane luluh.", NextSlideKey: "47"},
		{Key: "47", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Amin. Matur nuwun sanget wejanganipun nggih Bu.", NextSlideKey: "48"},

		// choice 8
		{
			Key: "48", Speaker: "Andi", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("neutral"), butejo("neutral")},
			Content: "(Badhe wangsul) Nggih pun Bu...",
			Choices: []choiceSeedData{
				{Text: "{Kula} mulih riyen.", NextSlideKey: "49a", MoodImpact: 0},
				{Text: "{Kula} {badhe} {wangsul}.", NextSlideKey: "49b", MoodImpact: 1},
				{Text: "{Kula} nuwun pamit.", NextSlideKey: "49c", MoodImpact: 1},
			},
			VocabKeys: []string{"kula", "badhe", "wangsul"},
		},

		// branch choice 8
		{Key: "49a", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("nervous"), butejo("angry")}, Content: "Walah, diwarahi ket mau kok bali 'Mulih' maneh. '{Wangsul}' Mas! Mulih iku kasar.", NextSlideKey: "50", VocabKeys: []string{"wangsul"}},
		{Key: "49b", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Ati-ati neng dalan.", NextSlideKey: "50"},
		{Key: "49c", Speaker: "Bu Tejo", BgImg: "bg/toko_batik.webp", Characters: []charData{andi("happy"), butejo("happy")}, Content: "Monggo-monggo. Ati-ati nggih Mas Ganteng. Salam kagem calonipun.", NextSlideKey: "50"},

		// closing
		{Key: "50", Speaker: "Narator", BgImg: "bg/toko_batik.webp", Content: "Andi medal saking toko kanthi eseman lega. Misi {madosi} oleh-oleh sampun rampung kanthi sukses.", NextSlideKey: "51", VocabKeys: []string{"madosi"}},
		{Key: "51", Speaker: "Narator", BgImg: "bg/toko_batik.webp", Content: "Pranyata leres ngendikanipun Pakdhe, ngadepi Bu Tejo mawon sampun kringeten, menapa malih Juragan Cengkeh asli.", NextSlideKey: "52"},
		{Key: "52", Speaker: "Narator", BgImg: "bg/toko_batik.webp", Content: "Nanging samenika, Andi sampun langkung siyap. Wancinipun budhal dhateng medan perang sejatine: Tulungagung."},
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
