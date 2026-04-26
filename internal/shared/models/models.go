package models

import "time"

type User struct {
	ID               uint64 `json:"id" gorm:"primaryKey"`
	Email            string `json:"email" gorm:"unique;not null"`
	PasswordHash     string `json:"-" gorm:"not null"`
	TwoFactorSecret  string `json:"-" gorm:"type:varchar(26)"`
	TwoFactorEnabled bool   `json:"two_factor_enabled" gorm:"default:false"`
	FirstSession     bool   `json:"first_session" gorm:"default:true"`
	FullProfile      bool   `json:"full_profile" gorm:"default:false"`
	EmailConfirmed   bool   `json:"email_confirmed" gorm:"default:false"`

	Token string `json:"token" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type File struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	FileID   string `json:"file_id" gorm:"unique;index"`
	Path     string `json:"-" gorm:"type:text;not null"`
	MimeType string `json:"mime_type" gorm:"not null"` // "image/jpeg",
	Size     int64  `json:"size" gorm:"not null"`      // en bytes}
	UserID   uint64 `json:"-" gorm:"not null"`

	Url string `json:"url" gorm:"-"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

type Organizer struct {
	ID     uint64 `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"not null"`
	Phone  string `json:"phone" gorm:"not null"`
	UserID uint64 `json:"user_id" gorm:"not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

type RaffleStatus struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique;not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Raffle struct {
	ID               uint64  `json:"id" gorm:"primaryKey"`
	Title            string  `json:"title" gorm:"not null"`
	Description      string  `json:"description" gorm:"type:text"`
	Price            float64 `json:"price" gorm:"not null"`
	ImageURL         string  `json:"image_url" gorm:"type:text"`
	Date             string  `json:"date" gorm:"type:date;not null"`
	MaxWinners       uint64  `json:"max_winners" gorm:"default:1"`
	QuantityTickets  uint64  `json:"quantity_tickets" gorm:"not null"`
	Slug             string  `json:"slug" gorm:"type:varchar(255);unique;not null"`
	OrganizerID      uint64  `json:"organizer_id" gorm:"not null"`
	RaffleStatusID   uint64  `json:"raffle_status_id" gorm:"not null;default:1"` // 1 = Activa, 2 = Cerrada, 3 = Sorteada, etc.
	TicketsAvailable uint64  `json:"tickets_available" gorm:"-"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Organizer    Organizer    `json:"-" gorm:"foreignKey:OrganizerID"`
	RaffleStatus RaffleStatus `json:"-" gorm:"foreignKey:RaffleStatusID"`
}

type Prize struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description" gorm:"type:text"`
	ImageURL    string `json:"image_url" gorm:"type:text"`
	OrganizerID uint64 `json:"organizer_id" gorm:"not null"`
	RaffleID    uint64 `json:"raffle_id" gorm:"-"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Organizer Organizer `json:"-" gorm:"foreignKey:OrganizerID"`
}

type RafflePrize struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	RaffleID uint64 `json:"raffle_id" gorm:"not null"`
	PrizeID  uint64 `json:"prize_id" gorm:"not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Raffle Raffle `json:"raffle" gorm:"foreignKey:RaffleID"`
	Prize  Prize  `json:"prize" gorm:"foreignKey:PrizeID"`
}

type TicketStatus struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique;not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Ticket struct {
	ID               uint64 `json:"id" gorm:"primaryKey"`
	Number           uint64 `json:"number" gorm:"not null"`
	ParticipantName  string `json:"participant_name" gorm:"type:varchar(255)"`
	ParticipantPhone string `json:"participant_phone" gorm:"type:varchar(20)"`
	RaffleID         uint64 `json:"raffle_id" gorm:"not null"`
	TicketStatusID   uint64 `json:"ticket_status_id" gorm:"not null;default:2"`
	Winner           bool   `json:"winner" gorm:"default:false"`
	OrganizerNumber  string `json:"organizer_number" gorm:"-"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Raffle       Raffle       `json:"-" gorm:"foreignKey:RaffleID"`
	TicketStatus TicketStatus `json:"ticket_status" gorm:"foreignKey:TicketStatusID"`
}

var Models = []any{
	&File{},
	&User{},
	&Organizer{},
	&Prize{},
	&RaffleStatus{},
	&Raffle{},
	&RafflePrize{},
	&TicketStatus{},
	&Ticket{},
}

var RaffleStatuses = []RaffleStatus{
	{ID: 1, Name: "Activa"},
	{ID: 2, Name: "Cerrada"},
	{ID: 3, Name: "Finalizada"},
}

var TicketStatuses = []TicketStatus{
	{ID: 1, Name: "Vendido"},
	{ID: 2, Name: "Disponible"},
	{ID: 3, Name: "Reservado"},
	{ID: 4, Name: "Anulado"},
}
