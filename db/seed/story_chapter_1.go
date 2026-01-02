package seed

import (
	"encoding/json"
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/internal/domain/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedChapter1(db *gorm.DB, vocabMap map[string]uuid.UUID) error {
	slog.Info("seeding chapter 1...")

	var chapter entity.Chapter
	err := db.Where("order_index = ?", 1).First(&chapter).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			chapter = entity.Chapter{
				Title:         "Ana Kabar Kaget!",
				Description:   "Andi kudu nekat budhal nang Tulungagung demi mjuangne tresnane.",
				CoverImageURL: "chapters/ch1_cover.webp",
				OrderIndex:    1,
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

	slidesData := []slideData{
		// intro
		{Key: "1", Speaker: "Narator", BgImg: "bg/warmindo.webp", Content: "Wanci {sonten} ing kutha Surabaya. Hawa panas taksih krasa, nanging ing satunggaling warung {alit}, swasana katingal ayem.", NextSlideKey: "2", VocabKeys: []string{"sonten", "alit"}},
		{Key: "2", Speaker: "Narator", BgImg: "bg/warmindo.webp", Content: "Warung {menika} namanipun 'Warmindo Andi'. Ingkang gadhah, satunggaling nom-noman ingkang grapyak lan remen guyon.", NextSlideKey: "3", VocabKeys: []string{"menika"}},

		{Key: "3", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("happy")}, Content: "(Nyelehake mangkok ing meja) Iki lho, Indomie telor kornet spesial! Mung gawe Cah Ayu sing paling manis sak Surabaya.", NextSlideKey: "4"},
		{Key: "4", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("happy")}, Content: "(Mesem) Matur nuwun lho, Mas Andi. Iki mesti enak banget.", NextSlideKey: "5"},
		{Key: "5", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("happy")}, Content: "Yo jelas no! Eh, tapi sik... koen kok ket {kalawau} meneng ae? Biasane cerewet nek ndelok drakor.", NextSlideKey: "6", VocabKeys: []string{"kalawau"}},
		{Key: "6", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "{Boten} kok, Mas. Aku mung lagi mikir {sakedhik}.", NextSlideKey: "7", VocabKeys: []string{"boten", "sakedhik"}},
		{Key: "7", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("neutral")}, Content: "Mikir opo? Mikir utang ta? Tenang, nek mangan neng kene gratis tis!", NextSlideKey: "8"},
		{Key: "8", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("happy")}, Content: "(Ngguyu cilik) Dudu kuwi, Mas. Iki babagan Bapak.", NextSlideKey: "9"},
		{Key: "9", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("happy")}, Content: "(Raine mulai serius) Bapakmu? Pak Broto sing juragan cengkeh iku? Lapo Bapakmu?", NextSlideKey: "10"},
		{Key: "10", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "Mau isuk Bapak {ngendikan} karo aku. Jarene... Bapak pengen {panggih} Mas Andi.", NextSlideKey: "11", VocabKeys: []string{"ngendikan", "panggih"}},
		{Key: "11", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("neutral")}, Content: "(Kaget, mripate mendelik) Hah?! Ketemu aku? Lapo? Aku ono salah ta?", NextSlideKey: "12"},
		{Key: "12", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("neutral")}, Content: "Ora ono sing salah, Mas. Bapak mung pengen kenalan. Jarene, 'Endi bocah lanang sing wani nyedaki anakku?' ngono.", NextSlideKey: "13"},
		{Key: "13", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("neutral")}, Content: "Waduh... mati aku. Koen eruh dewe Bapakmu kaya opo.", NextSlideKey: "14"},
		{Key: "14", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("happy")}, Content: "Aja ngono ta, Mas. Bapak ki asline apikan kok, mung trampil wae.", NextSlideKey: "15"},
		{Key: "15", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("happy")}, Content: "Trampil jaremu... wingi onok maling ayam ae meh dipenthung tongkat.", NextSlideKey: "16"},
		{Key: "16", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "Pokoke minggu ngarep Mas Andi kudu {sowan} nang Tulungagung. Bapak wis nunggu.", NextSlideKey: "17", VocabKeys: []string{"sowan"}},

		// choice 1
		{
			Key: "17", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("neutral")},
			Content: "(Garuk-garuk sirah sing ora gatel) Duh, piye iki...",
			Choices: []choiceSeedData{
				{Text: "Waduh Dek, aku durung siyap mental! Iso semaput aku pas salaman.", NextSlideKey: "18a", MoodImpact: -1},
				{Text: "Oke, sapa wedi! Bonek wani perih! Pak Broto sapa?", NextSlideKey: "18b", MoodImpact: 0},
			},
		},

		// branch choice 1
		{Key: "18a", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("worried")}, Content: "Mas Andi... jare sampeyan serius karo aku? Mosok ketemu Bapak wae wedi?", NextSlideKey: "19a"},
		{Key: "19a", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("worried")}, Content: "Gak ngono... tapi Bapakmu iku lho, aurane kaya Singo Barong.", NextSlideKey: "20a"},
		{Key: "20a", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("happy")}, Content: "Tenang wae, onok aku kok. Mas Andi sing penting niate apik.", NextSlideKey: "21"},

		{Key: "18b", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("happy")}, Content: "(Mesem ngenyek) Ooh, ngono ya? Yowes, ngko tak kandhakne Bapak nek Mas Andi nantang.", NextSlideKey: "19b"},
		{Key: "19b", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("happy")}, Content: "(Langsung ciut) Eh, ojo! Ojo dikandhakne 'nantang', Rek. Maksudku, aku wani sowan kanthi sopan.", NextSlideKey: "20b"},
		{Key: "20b", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("happy")}, Content: "Nah, ngono lho. Sing penting sopan.", NextSlideKey: "21"},

		// merge path
		{Key: "21", Speaker: "Narator", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("happy")}, Content: "Sanajan lesanipun matur wantun, nanging manahipun Andi tetep keraos deg-degan.", NextSlideKey: "22"},
		{Key: "22", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("happy")}, Content: "Tapi onok masalah siji maneh, Sekar.", NextSlideKey: "23"},
		{Key: "23", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "Opo maneh, Mas?", NextSlideKey: "24"},
		{Key: "24", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("neutral")}, Content: "Aku iki gak iso Boso Krama! Koen eruh dewe cangkemku iki isine misuh tok.", NextSlideKey: "25"},
		{Key: "25", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("neutral")}, Content: "Mosok pas ketemu Bapakmu aku ngomong, 'Pak, koen sehat ta?' Lak iso diusir aku!", NextSlideKey: "26"},
		{Key: "26", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("neutral")}, Content: "Iya ya... Bapak kuwi njawani banget, Mas. Mas Andi kudu nggawe Krama Inggil nek {ngendikan}.", NextSlideKey: "27", VocabKeys: []string{"ngendikan"}},
		{Key: "27", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "Nah iku! Aku mbien pelajaran Boso Jowo turu tok. Jajal tes aku.", NextSlideKey: "28"},

		// choice 2
		{
			Key: "28", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("happy")},
			Content: "Yowes, coba tak tes ya. Boso Kramane 'Mangan' nek kanggo Bapak iku opo?",
			Choices: []choiceSeedData{
				{Text: "Nedha", NextSlideKey: "29a", MoodImpact: 0},
				{Text: "Dhahar", NextSlideKey: "29b", MoodImpact: 1},
				{Text: "Badhog", NextSlideKey: "29c", MoodImpact: -1},
			},
		},

		// branch choice 2
		{Key: "29a", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "Kurang pas, Mas. '{Nedha}' iku kanggo awake dhewe (Krama Lugu). Nek kanggo Bapak, kudune '{Dhahar}'.", NextSlideKey: "30a", VocabKeys: []string{"nedha", "dhahar"}},
		{Key: "30a", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "Woalah ngono ta. Meh bener kan tapi.", NextSlideKey: "31"},

		{Key: "29b", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("happy")}, Content: "Wah, bener! Pinter Mas Andi. Kuwi tembung sing pas kanggo Bapak.", NextSlideKey: "30b"},
		{Key: "30b", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("happy")}, Content: "(Bangga) Lha iyo, sopo disek pacare.", NextSlideKey: "31"},

		{Key: "29c", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("angry")}, Content: "(Mendelik) Heh! Aja ngawur! Kuwi kasar banget, Mas! Iso dicoret saka KK kowe engko.", NextSlideKey: "30c"},
		{Key: "30c", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("angry")}, Content: "Waduh, iyo lali. Sepurane.", NextSlideKey: "31"},

		// merge path
		{Key: "31", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("neutral")}, Content: "Duh, angel men. Siji kata wae aku wis bingung, opo maneh sak kalimat.", NextSlideKey: "32"},
		{Key: "32", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("nervous"), sekar("neutral")}, Content: "Mas Andi kudu sinau. Wektune isih seminggu.", NextSlideKey: "33"},
		{Key: "33", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "Aku kudu njaluk tulung sopo ya? Koen iso ngajari aku ta?", NextSlideKey: "34"},
		{Key: "34", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "Aku iso {sakedhik}, Mas. Tapi luwih apik nek Mas Andi sinau karo wong sing luwih sepuh.", NextSlideKey: "35", VocabKeys: []string{"sakedhik"}},
		{Key: "35", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "Soale unggah-ungguh kuwi ora mung basa, tapi uga sikap.", NextSlideKey: "36"},
		{Key: "36", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("neutral"), sekar("neutral")}, Content: "(_Mikir_) Wong sepuh sing santai tapi pinter...", NextSlideKey: "37"},
		{Key: "37", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("neutral")}, Content: "Ahh! Pakdhe Joyo!", NextSlideKey: "38"},
		{Key: "38", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("happy")}, Content: "Lha, Pakdhe Joyo kan pinter banget Basa Jawi. Panjenengan {sowan} rono wae, Mas.", NextSlideKey: "39", VocabKeys: []string{"sowan"}},
		{Key: "39", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("happy")}, Content: "Bener! Pakdhe Joyo iku wis koyo guruku dewe. Yowes, sesuk isuk aku tak mrono.", NextSlideKey: "40"},
		{Key: "40", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("happy")}, Content: "Sip. Mas Andi kudu semangat ya. Iki demi masa depan lho.", NextSlideKey: "41"},
		{Key: "41", Speaker: "Andi", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("happy")}, Content: "Demi awakmu, Sekar. Opo sih sing ora tak lakoni.", NextSlideKey: "42"},
		{Key: "42", Speaker: "Sekar", BgImg: "bg/warmindo.webp", Characters: []charData{andi("happy"), sekar("happy")}, Content: "(Isin-isin) Halah, gombal.", NextSlideKey: "43"},
		{Key: "43", Speaker: "Narator", BgImg: "bg/warmindo.webp", Content: "Mekatenlah wiwitanipun perjuangan Andi. Saking warung mie, tumuju manahipun calon mertua.", NextSlideKey: "44"},
		{Key: "44", Speaker: "Narator", BgImg: "bg/warmindo.webp", Content: "Punapa Andi {badhe} kasil sinau tata krama saking Pakdhe Joyo? Entosi cariyos salajengipun.", VocabKeys: []string{"badhe"}},
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
