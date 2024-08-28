package ad

import (
	"context"
	"slices"
	"time"
)

type UserRepo interface {
	GetByID(context.Context, UserID)
}

func NewAdService(userRepo) *adServ
func (s *adServ) CreateAdDraft(ctx, userID /*...*/) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if !user.IsBanned() {
		return nil, nil
	}
	/*...*/
}

func Archive(ctx, adID int) error {
	ad, err := db.Get(ctx, adID)
	// Архивируем
	ad.ArchivedAt = time.Now()

	err = db.Save(ctx, ad)

	return err
}

func Archive(ctx, adID int) error {
	ad, err := db.Get(ctx, adID)
	// Архивируем
	if !a.Status.Can(Archived) {
		return ErrArchive
	}
	a.Status = Archived
	ad.ArchivedAt = time.Now()

	err = db.Save(ctx, ad)

	return nil
}

type AdID int

func NewAdID(id int) (AdID, error) { /*...*/ }

const (
	Draft Status = iota + 1
	Active
	Archived
)

type Status int

func (s Status) CanBe() bool {
	switch s {
	case Draft:
		return []Status{Active}
		// Чтобы учесть все варианты Status
		// линтер - exhaustive
	}
	//...
}

func (s Status) Can(to Status) bool {
	return slices.Contains(s.CanBe(), to)
}

type Status int
const Draft Status = 1

ad := &Ad{
Status: Draft,
Attrs: nil,
}

// panic: assignment to entry in nil map
ad.Attrs["key"] = "value"
