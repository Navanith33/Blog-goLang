package models
type Blog struct {
	ID        uint        `gorm:"primaryKey;autoIncrement"`
	Title     string      `gorm:"not null;size:255"`
	Content   string      `gorm:"not null"`
	UserId    uint
	User      User        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
  }
  type User struct {
	ID           uint     `gorm:"primaryKey;autoIncrement"`
	Email        string   `gorm:"unique;size:255;not null"`  
	Password    string   `gorm:"not null"`
	Role        string
	Blogs        []Blog
  }