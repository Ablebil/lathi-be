package seed

import (
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BadgeSeeder struct{}

func (s *BadgeSeeder) Run(db *gorm.DB) error {
	slog.Info("seeding badge domain...")

	badges := []entity.Badge{
		{
			Code:        "ch1_completion",
			Name:        "Satriya Anyar",
			Description: "Langkah kapisan wis dilakoni. Wani miwiti, wani mungkasi.",
			IconURL:     "badges/satriya_anyar.webp",
		},
		{
			Code:        "vocab_collector_1",
			Name:        "Widya Tembung",
			Description: "Bausastra ing sirah saya akeh. Bekal wicara supaya ora kleru.",
			IconURL:     "badges/widya_tembung.webp",
		},
		{
			Code:        "perfect_heart",
			Name:        "Susila Anoraga",
			Description: "Alus ing rasa, manis ing basa. Tetep santun tanpa gawe gelane liyan.",
			IconURL:     "badges/susila_anoraga.webp",
		},
		{
			Code:        "all_chapters_completion",
			Name:        "Wiyata Paripurna",
			Description: "Lulus kabeh pacoban. Sampun pantes lungguh ing ngarep.",
			IconURL:     "badges/wiyata_paripurna.webp",
		},
	}

	for _, b := range badges {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "description", "icon_url"}),
		}).Create(&b).Error; err != nil {
			slog.Error("failed to seed badge", "code", b.Code, "error", err)
			return err
		}
	}

	slog.Info("badge seeding completed successfully")
	return nil
}
