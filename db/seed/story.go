package seed

import (
	"encoding/json"
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/internal/domain/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorySeeder struct{}

type charData struct {
	Name string
	Img  string
}

type choiceSeedData struct {
	Text         string
	NextSlideKey string
	MoodImpact   int
}

type slideData struct {
	Key          string
	Speaker      string
	Content      string
	BgImg        string
	Characters   []charData
	NextSlideKey string
	Choices      []choiceSeedData
	VocabKeys    []string
}

func (s *StorySeeder) Run(db *gorm.DB) error {
	slog.Info("seeding story domain...")

	vocabs := []entity.Dictionary{
		// ch 1
		{WordKrama: "sonten", WordNgoko: "sore", WordIndo: "sore"},
		{WordKrama: "alit", WordNgoko: "cilik", WordIndo: "kecil"},
		{WordKrama: "menika", WordNgoko: "iki/kuwi", WordIndo: "ini/itu"},
		{WordKrama: "kalawau", WordNgoko: "mau", WordIndo: "tadi"},
		{WordKrama: "ngendikan", WordNgoko: "ngomong", WordIndo: "berbicara"},
		{WordKrama: "panggih", WordNgoko: "ketemu", WordIndo: "bertemu"},
		{WordKrama: "sowan", WordNgoko: "teka/mertamu", WordIndo: "berkunjung (hormat)"},
		{WordKrama: "boten", WordNgoko: "ora", WordIndo: "tidak"},
		{WordKrama: "sakedhik", WordNgoko: "sithik", WordIndo: "sedikit"},
		{WordKrama: "dhahar", WordNgoko: "mangan", WordIndo: "makan (hormat)"},
		{WordKrama: "nedha", WordNgoko: "mangan", WordIndo: "makan (umum)"},
		{WordKrama: "badhe", WordNgoko: "arep", WordIndo: "akan"},

		// ch 2
		{WordKrama: "enjang", WordNgoko: "esuk", WordIndo: "pagi"},
		{WordKrama: "saweg", WordNgoko: "lagi", WordIndo: "sedang"},
		{WordKrama: "ngunjuk", WordNgoko: "ngombe", WordIndo: "minum (hormat)"},
		{WordKrama: "pripun", WordNgoko: "piye", WordIndo: "bagaimana"},
		{WordKrama: "wonten", WordNgoko: "ana", WordIndo: "ada"},
		{WordKrama: "kula", WordNgoko: "aku", WordIndo: "saya"},
		{WordKrama: "dalem", WordNgoko: "aku", WordIndo: "saya (sangat hormat)"},
		{WordKrama: "panjenengan", WordNgoko: "kowe", WordIndo: "anda (hormat)"},
		{WordKrama: "sampeyan", WordNgoko: "kowe", WordIndo: "kamu (umum)"},
		{WordKrama: "saking", WordNgoko: "saka", WordIndo: "dari"},
		{WordKrama: "lare", WordNgoko: "bocah", WordIndo: "anak/orang"},
		{WordKrama: "griya", WordNgoko: "omah", WordIndo: "rumah"},
		{WordKrama: "wangsul", WordNgoko: "bali/mulih", WordIndo: "pulang"},

		// ch 3
		{WordKrama: "tumbas", WordNgoko: "tuku", WordIndo: "beli (umum)"},
		{WordKrama: "mundhut", WordNgoko: "tuku", WordIndo: "membeli (hormat)"},
		{WordKrama: "awis", WordNgoko: "larang", WordIndo: "mahal"},
		{WordKrama: "arta", WordNgoko: "duwit", WordIndo: "uang"},
		{WordKrama: "regi", WordNgoko: "rego", WordIndo: "harga"},
		{WordKrama: "paring", WordNgoko: "wenehi", WordIndo: "memberi"},
		{WordKrama: "kersa", WordNgoko: "gelem", WordIndo: "mau/bersedia"},
		{WordKrama: "remen", WordNgoko: "seneng", WordIndo: "suka/senang"},
		{WordKrama: "kagungan", WordNgoko: "duwe", WordIndo: "milik/mempunyai"},
		{WordKrama: "setunggal", WordNgoko: "siji", WordIndo: "satu"},
		{WordKrama: "kalih", WordNgoko: "loro", WordIndo: "dua"},
		{WordKrama: "tiga", WordNgoko: "telu", WordIndo: "tiga"},
		{WordKrama: "sekawan", WordNgoko: "papat", WordIndo: "empat"},
		{WordKrama: "gangsal", WordNgoko: "lima", WordIndo: "lima"},
		{WordKrama: "dasa", WordNgoko: "sepuluh", WordIndo: "sepuluh"},
		{WordKrama: "seket", WordNgoko: "seket", WordIndo: "lima puluh"},
		{WordKrama: "atus", WordNgoko: "atus", WordIndo: "ratus"},
		{WordKrama: "ewu", WordNgoko: "ewu", WordIndo: "ribu"},
		{WordKrama: "pinten", WordNgoko: "piro", WordIndo: "berapa"},
		{WordKrama: "niku", WordNgoko: "kuwi", WordIndo: "itu (umum)"},
		{WordKrama: "pundi", WordNgoko: "endi", WordIndo: "mana"},
		{WordKrama: "madosi", WordNgoko: "golek", WordIndo: "mencari"},

		// ch 4
		{WordKrama: "kula nuwun", WordNgoko: "permisi", WordIndo: "permisi (saat masuk rumah)"},
		{WordKrama: "mlebet", WordNgoko: "mlebu", WordIndo: "masuk"},
		{WordKrama: "lenggah", WordNgoko: "lungguh", WordIndo: "duduk"},
		{WordKrama: "sugeng", WordNgoko: "slamet", WordIndo: "selamat"},
		{WordKrama: "leres", WordNgoko: "bener", WordIndo: "benar"},
		{WordKrama: "sadeyan", WordNgoko: "dodolan", WordIndo: "berjualan"},
		{WordKrama: "cekap", WordNgoko: "cukup", WordIndo: "cukup"},
		{WordKrama: "ajrih", WordNgoko: "wedi", WordIndo: "takut/segan"},
		{WordKrama: "manah", WordNgoko: "ati", WordIndo: "hati/perasaan"},
		{WordKrama: "estu", WordNgoko: "tenan", WordIndo: "sungguh/benar-benar"},
		{WordKrama: "nyuwun", WordNgoko: "njaluk", WordIndo: "meminta"},
		{WordKrama: "ngrantos", WordNgoko: "ngenteni", WordIndo: "menunggu"},
		{WordKrama: "asrep", WordNgoko: "adhem", WordIndo: "dingin"},
		{WordKrama: "dhateng", WordNgoko: "menyang", WordIndo: "ke/kepada"},
		{WordKrama: "wancinipun", WordNgoko: "wayahe", WordIndo: "waktunya"},
		{WordKrama: "ngarsanipun", WordNgoko: "ngarepe", WordIndo: "di hadapan"},
		{WordKrama: "dugi", WordNgoko: "teka", WordIndo: "tiba/sampai"},
	}

	vocabMap := make(map[string]uuid.UUID)
	for _, v := range vocabs {
		var dict entity.Dictionary
		err := db.Where("word_krama = ?", v.WordKrama).First(&dict).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&v).Error; err != nil {
					return err
				}
				dict = v
			} else {
				return err
			}
		}
		vocabMap[v.WordKrama] = dict.ID
	}

	// execute chapter seeders
	if err := seedChapter1(db, vocabMap); err != nil {
		slog.Error("failed to seed chapter 1", "error", err)
		return err
	}

	if err := seedChapter2(db, vocabMap); err != nil {
		slog.Error("failed to seed chapter 2", "error", err)
		return err
	}

	if err := seedChapter3(db, vocabMap); err != nil {
		slog.Error("failed to seed chapter 3", "error", err)
		return err
	}

	if err := seedChapter4(db, vocabMap); err != nil {
		slog.Error("failed to seed chapter 4", "error", err)
		return err
	}

	slog.Info("story seeding completed successfully")
	return nil
}

func makeChoicesWithRealIDs(opts []choiceSeedData, realIDs map[string]uuid.UUID) types.JSONB {
	type choiceJSON struct {
		Text        string    `json:"text"`
		NextSlideID uuid.UUID `json:"next_slide_id"`
		MoodImpact  int       `json:"mood_impact"`
	}

	res := make([]choiceJSON, len(opts))
	for i, o := range opts {
		nextID := realIDs[o.NextSlideKey]
		res[i] = choiceJSON{
			Text:        o.Text,
			NextSlideID: nextID,
			MoodImpact:  o.MoodImpact,
		}
	}

	b, _ := json.Marshal(res)
	return types.JSONB(b)
}
