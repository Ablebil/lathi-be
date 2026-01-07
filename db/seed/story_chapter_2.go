package seed

import (
	"encoding/json"
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/internal/domain/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedChapter2(db *gorm.DB, vocabMap map[string]uuid.UUID) error {
	slog.Info("seeding chapter 2...")

	var chapter entity.Chapter
	err := db.Where("order_index = ?", 2).First(&chapter).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			chapter = entity.Chapter{
				Title:         "Sinau Dadi Priyayi",
				Description:   "Andi meguru tata krama nang Pakdhe Joyo ben ora ngisin-ngisini.",
				CoverImageURL: "chapters/ch2_cover.webp",
				OrderIndex:    2,
			}
			if err := db.Create(&chapter).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	andi := func(exp string) charData { return charData{Name: "Andi", Img: "chars/andi_" + exp + ".webp"} }
	pakdhe := func(exp string) charData { return charData{Name: "Pakdhe Joyo", Img: "chars/pakdhe_" + exp + ".webp"} }

	slidesData := []slideData{
		// intro
		{Key: "1", Speaker: "Narator", BgImg: "bg/teras_joglo.webp", Content: "{Enjang} menika, srengenge katingal sumunar padhang. Manuk perkutut manggung saut-sautan ing teras omah Joglo.", NextSlideKey: "2", VocabKeys: []string{"enjang"}},
		{Key: "2", Speaker: "Narator", BgImg: "bg/teras_joglo.webp", Content: "Andi sampun dumugi ing dalemipun Pakdhe Joyo, sesepuh ingkang badhe dipunsuwuni pirsa.", NextSlideKey: "3"},

		// choice 1
		{
			Key: "3", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("neutral")},
			Content: "(Mlebu pager, ndelok Pakdhe lagi lungguh maca koran) (_Wah, Pakdhe ketoke santai banget. Kudu nyapa sing sopan iki._)",
			Choices: []choiceSeedData{
				{Text: "Sugeng {enjang} Pakdhe, {saweg} {ngunjuk} kopi niki?", NextSlideKey: "4a", MoodImpact: 1},
				{Text: "{Pripun} kabare Dhe? Sehat?", NextSlideKey: "4b", MoodImpact: 0},
				{Text: "Wooy Dhe! Lagi ngopi ta?", NextSlideKey: "4c", MoodImpact: -1},
			},
			VocabKeys: []string{"enjang", "saweg", "ngunjuk", "pripun"},
		},

		// branch choice 1
		{Key: "4a", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "(Noleh kaget seneng) Weh, tumben kowe boso, Le. Iyo, kene lungguh. Kadingaren men.", NextSlideKey: "5a"},
		{Key: "5a", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Inggih Dhe, nuwun sewu.", NextSlideKey: "6"},

		{Key: "4b", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("neutral")}, Content: "Apik-apik wae Le. Tumben isuk-isuk mrene.", NextSlideKey: "5b"},
		{Key: "5b", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("neutral")}, Content: "Inggih Dhe, {wonten} perlu.", NextSlideKey: "6", VocabKeys: []string{"wonten"}},

		{Key: "4c", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("angry")}, Content: "(Gedheg-gedheg) Hush! Bengok-bengok kaya neng alas! Iki omah Le, udu terminal.", NextSlideKey: "5c"},
		{Key: "5c", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("angry")}, Content: "Waduh, sepurane Dhe. Kabotan cangkem.", NextSlideKey: "6"},

		// merge path
		{Key: "6", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("neutral")}, Content: "Dhe, aku butuh bantuanmu. Minggu ngarep aku dikongkon nemoni Bapake Sekar, Pak Broto Juragan Cengkeh.", NextSlideKey: "7"},
		{Key: "7", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("neutral")}, Content: "Walah Le... Pak Broto iku wong 'Priyayi Sepuh'. Kowe kudu ati-ati, tata kramane dijogo tenan.", NextSlideKey: "8"},
		{Key: "8", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("neutral")}, Content: "Lha nggih niku! Kula gak iso Boso Krama, Dhe. Tulung warahono kula.", NextSlideKey: "9"},
		{Key: "9", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "Yowes. Tak warahi sing dhasar wae. Sing penting 'Unggah-Ungguh Basa' lan 'Solah Bawa'.", NextSlideKey: "10"},

		// choice 2
		{
			Key: "10", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")},
			Content: "Sepisan, nek nyebut awakmu dewe neng ngarepe Pak Broto, kowe nggawe tembung apa?",
			Choices: []choiceSeedData{
				{Text: "{Kula}", NextSlideKey: "11a", MoodImpact: 1},
				{Text: "{Dalem}", NextSlideKey: "11b", MoodImpact: 0},
				{Text: "Aku", NextSlideKey: "11c", MoodImpact: -1},
			},
			VocabKeys: []string{"kula", "dalem"},
		},

		// branch choice 2
		{Key: "11a", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Bener. '{Kula}' iku standar sopan gawe awake dhewe.", NextSlideKey: "12", VocabKeys: []string{"kula"}},
		{Key: "11b", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("neutral")}, Content: "'{Dalem}' iku apik banget, tapi kadhang krasa kaya abdi dalem. Kanggo wiwitan, '{Kula}' wae wis cukup.", NextSlideKey: "12", VocabKeys: []string{"dalem", "kula"}},
		{Key: "11c", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("angry")}, Content: "Salah! 'Aku' iku Ngoko. Kasar nek gawe ngomong karo calon mertua.", NextSlideKey: "12"},

		// choice 3
		{
			Key: "12", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")},
			Content: "Lha nek nyebut Pak Broto? Kowe nyeluk piye?",
			Choices: []choiceSeedData{
				{Text: "{Panjenengan}", NextSlideKey: "13a", MoodImpact: 1},
				{Text: "{Sampeyan}", NextSlideKey: "13b", MoodImpact: 0},
				{Text: "Kowe", NextSlideKey: "13c", MoodImpact: -1},
			},
			VocabKeys: []string{"panjenengan", "sampeyan"},
		},

		// branch choice 3
		{Key: "13a", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Pinter. '{Panjenengan}' iku paling alus. Pak Broto mesthi seneng nek diceluk ngono.", NextSlideKey: "14", VocabKeys: []string{"panjenengan"}},
		{Key: "13b", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("neutral")}, Content: "Aja '{Sampeyan}', Le. Iku isih rodok kasar gawe priyayi sepuh. Iku gawe kancan utawa wong sing sak level.", NextSlideKey: "14", VocabKeys: []string{"sampeyan"}},
		{Key: "13c", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("angry")}, Content: "Weh, ojo pisan-pisan nggawe 'Kowe'. Iku ngoko kasar! Iso langsung diusir kowe.", NextSlideKey: "14"},

		// merge path
		{Key: "14", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "Wis, dhasare iku sik. Saiki kowe kudu latihan saben dina. Aja mung teori tok.", NextSlideKey: "15"},
		{Key: "15", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Siap Dhe. Saben sore aku bakal mrene maneh.", NextSlideKey: "16"},

		{Key: "16", Speaker: "Narator", BgImg: "bg/teras_joglo.webp", Content: "Dinten-dinten salajengipun, Andi sregep sowan dhateng dalemipun Pakdhe Joyo.", NextSlideKey: "17"},
		{Key: "17", Speaker: "Narator", BgImg: "bg/teras_joglo.webp", Content: "Saben sonten, Andi sinau cara lenggah, cara mlampah, lan cara matur ingkang leres.", NextSlideKey: "18"},

		{Key: "18", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "(Latihan mlaku mbungkuk) Ngeten nggih Dhe? Nuwun sewu...", NextSlideKey: "19"},
		{Key: "19", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "Kurang mbungkuk sithik. Aja kaku-kaku, sing luwes.", NextSlideKey: "20"},
		{Key: "20", Speaker: "Narator", BgImg: "bg/teras_joglo.webp", Content: "Ngantos dumugi dinten H-1, Pakdhe Joyo nguji asil pasinaonipun Andi.", NextSlideKey: "21"},

		{Key: "21", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "Nah, sesuk kowe wis arep mangkat nang Tulungagung. Saiki tak tes terakhir.", NextSlideKey: "22"},
		{Key: "22", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "Anggep aku Pak Broto. Aku nawani kowe mangan. 'Monggo {dhahar} riyen, Nak Andi'.", NextSlideKey: "23", VocabKeys: []string{"dhahar"}},

		// choice 4
		{
			Key: "23", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")},
			Content: "(Mikir jawaban sing pas)",
			Choices: []choiceSeedData{
				{Text: "Inggih Pak, matur nuwun. Kula {nedha} sakmenika.", NextSlideKey: "24a", MoodImpact: 0},
				{Text: "Inggih Pak, matur nuwun. Sampun repot-repot.", NextSlideKey: "24b", MoodImpact: 1},
				{Text: "Inggih Pak, matur nuwun. Kula badhe {dhahar}.", NextSlideKey: "24c", MoodImpact: -1},
			},
			VocabKeys: []string{"nedha", "dhahar"},
		},

		// branch choice 4
		{Key: "24a", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("neutral")}, Content: "'{Nedha}' iku bener krama, tapi nek neng ngarepe priyayi, luwih apik nggawe basa sing luwih umum dhisik.", NextSlideKey: "25", VocabKeys: []string{"nedha"}},
		{Key: "24b", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Sip. Jawaban sing aman lan sopan. Ora kakehan basa-basi.", NextSlideKey: "25"},
		{Key: "24c", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("angry")}, Content: "Heh! '{Dhahar}' iku kanggone wong liya sing dihormati! Nek kanggo awake dhewe iku '{Nedha}' utawa 'Maem'. Aja walik-walik!", NextSlideKey: "25", VocabKeys: []string{"dhahar", "nedha"}},

		// choice 5
		{
			Key: "25", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")},
			Content: "Saiki babagan pamitan. Kowe wis mari bertamu, arep mulih. Ngomong piye?",
			Choices: []choiceSeedData{
				{Text: "Nuwun sewu Pak, kula badhe {wangsul}.", NextSlideKey: "26a", MoodImpact: 0},
				{Text: "Nuwun sewu Pak, kula nyuwun pamit.", NextSlideKey: "26b", MoodImpact: 1},
				{Text: "Pak, aku mulih dhisik ya.", NextSlideKey: "26c", MoodImpact: -1},
			},
			VocabKeys: []string{"wangsul"},
		},

		// branch choice 5
		{Key: "26a", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("neutral")}, Content: "'{Wangsul}' iku bener, nanging 'Nuwun Pamit' luwih trep lan luwih alus kanggo tamu.", NextSlideKey: "27", VocabKeys: []string{"wangsul"}},
		{Key: "26b", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Joss! Iku kalimat pamungkas sing apik.", NextSlideKey: "27"},
		{Key: "26c", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("angry")}, Content: "Walah Le... 'Mulih' iku ngoko. Nek ngomong ngono, kowe ora bakal diundang meneh.", NextSlideKey: "27"},

		// choice 6
		{
			Key: "27", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")},
			Content: "Tes terakhir. Kowe ditakoni Pak Broto: 'Kowe asline wong endi?'",
			Choices: []choiceSeedData{
				{Text: "Kula {saking} Surabaya, Pak.", NextSlideKey: "28a", MoodImpact: 1},
				{Text: "{Dalem} asli {lare} Suroboyo, Pak.", NextSlideKey: "28b", MoodImpact: 1},
				{Text: "Omahku Surabaya, Pak.", NextSlideKey: "28c", MoodImpact: -1},
			},
			VocabKeys: []string{"saking", "dalem", "lare"},
		},

		// branch choice 6
		{Key: "28a", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("neutral")}, Content: "Cukup. Jawaban sing jelas lan ora neko-neko.", NextSlideKey: "29", VocabKeys: []string{"saking"}},
		{Key: "28b", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Wah, jan alus tenan. Pak Broto mesthi kaget nek krungu iki.", NextSlideKey: "29", VocabKeys: []string{"dalem", "lare"}},
		{Key: "28c", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("angry")}, Content: "Salah maneh. Aja nggawe 'Omahku'. Kudune '{Griya} Kula' utawa 'Kula {saking}...'.", NextSlideKey: "29", VocabKeys: []string{"griya", "saking"}},

		// merge path
		{Key: "29", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "Wis, secara teori kowe wis lulus. Tapi kowe butuh praktek nyata.", NextSlideKey: "30"},
		{Key: "30", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "Praktek pripun Dhe?", NextSlideKey: "31"},
		{Key: "31", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "Kowe kudu nggawa oleh-oleh kanggo Pak Broto. Tuku o Batik utawa Gethuk.", NextSlideKey: "32"},
		{Key: "32", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("teaching")}, Content: "Lungo o nang Toko Bu Tejo neng pasar. Jajal kowe nawar rego nggawe Basa Krama.", NextSlideKey: "33"},
		{Key: "33", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("neutral")}, Content: "Waduh, Bu Tejo sing galak kae ta?", NextSlideKey: "34"},
		{Key: "34", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("nervous"), pakdhe("happy")}, Content: "Iyo. Nek kowe iso ngadepi Bu Tejo tanpa diseneni, berarti kowe wis siap ngadepi Pak Broto.", NextSlideKey: "35"},

		// choice 7
		{
			Key: "35", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")},
			Content: "Siap Dhe! Tantangan ditampa.",
			Choices: []choiceSeedData{
				{Text: "Nggih Dhe, kula badhe nyobi.", NextSlideKey: "36a", MoodImpact: 1},
				{Text: "Oke Dhe, tak budhal saiki.", NextSlideKey: "36b", MoodImpact: 0},
				{Text: "Inggih Dhe, sendika dawuh.", NextSlideKey: "36c", MoodImpact: 1},
			},
		},

		// branch choice 7
		{Key: "36a", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Apik. Niat sing lurus iku kunci.", NextSlideKey: "37"},
		{Key: "36b", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("neutral"), pakdhe("neutral")}, Content: "Lho, kok bali ngoko maneh? Sing konsisten Le, sanajan karo Pakdhe.", NextSlideKey: "37"},
		{Key: "36c", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Weh weh, kok kaya abdi dalem Keraton. Tapi yo rapopo, semangatmu apik.", NextSlideKey: "37"},

		// closing
		{Key: "37", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("teaching")}, Content: "Ati-ati, Bu Tejo kae wonge rada cerewet. Sing sabar.", NextSlideKey: "38"},
		{Key: "38", Speaker: "Andi", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("teaching")}, Content: "Beres Dhe. Kula nuwun pamit. Assalamualaikum.", NextSlideKey: "39"},
		{Key: "39", Speaker: "Pakdhe Joyo", BgImg: "bg/teras_joglo.webp", Characters: []charData{andi("happy"), pakdhe("happy")}, Content: "Waalaikumsalam.", NextSlideKey: "40"},
		{Key: "40", Speaker: "Narator", BgImg: "bg/teras_joglo.webp", Content: "Kanthi bekal ilmu ingkang sampun dipunsinaoni pirang-pirang dinten, Andi tumuju pasar.", NextSlideKey: "41"},
		{Key: "41", Speaker: "Narator", BgImg: "bg/teras_joglo.webp", Content: "Punapa Andi badhe kasil ngluluhaken manahipun Bu Tejo kanthi basa kramanipun? Entosi cariyos salajengipun."},
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
