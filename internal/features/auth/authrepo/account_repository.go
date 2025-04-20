package authrepo

import (
	"context"
	"errors"
	"fmt"
	"harmony/internal/couchdb"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
)

var ErrConflict = couchdb.ErrConflict

type accountEmailDoc struct {
	authdomain.AccountID
}

type accountPasswordDoc struct {
	ID           authdomain.AccountID
	PasswordHash []byte
}

type AccountRepository struct {
	couchdb.Connection
}

func (r AccountRepository) accDocId(id authdomain.AccountID) string {
	return fmt.Sprintf("auth:account:%s", id)
}

func (r AccountRepository) addrDocId(addr string) string {
	return fmt.Sprintf("auth:account:email:%s", addr)
}

func (r AccountRepository) accEmailDocID(acc authdomain.Account) string {
	return r.addrDocId(acc.Email.String())
}

func passwordDocId(id authdomain.AccountID) string {
	return fmt.Sprintf("auth:accunt:%s:password", id)
}

func (r AccountRepository) insertAccountDoc(ctx context.Context, acc authdomain.Account) error {
	_, err := r.Connection.Insert(ctx, r.accDocId(acc.ID), acc)
	return err
}

func (r AccountRepository) insertEmailDoc(
	ctx context.Context,
	acc auth.AccountUseCaseResult,
) error {
	doc := couchdb.DocumentWithEvents[accountEmailDoc]{
		Document: accountEmailDoc{acc.Entity.ID},
		Events:   acc.Events,
	}
	_, err := r.Connection.Insert(ctx,
		r.accEmailDocID(acc.Entity.Account),
		doc,
	)
	return err
}

func (r AccountRepository) insertPasswordDoc(
	ctx context.Context,
	acc authdomain.PasswordAuthentication,
) error {
	doc := accountPasswordDoc{
		acc.ID,
		acc.PasswordHash.UnsecureRead(),
	}
	_, err := r.Connection.Insert(ctx, passwordDocId(acc.ID), doc)
	return err
}

func (r AccountRepository) Insert(
	ctx context.Context,
	acc auth.AccountUseCaseResult,
) error {
	err := r.insertAccountDoc(ctx, acc.Entity.Account)
	if err == nil {
		err = r.insertPasswordDoc(ctx, acc.Entity)
	}
	if err == nil {
		err = r.insertEmailDoc(ctx, acc)
	}
	return err
}

func (r AccountRepository) Get(
	ctx context.Context,
	id authdomain.AccountID,
) (res authdomain.Account, err error) {
	_, err = r.Connection.Get(r.accDocId(id), &res)
	return
}

func (r AccountRepository) FindByEmail(ctx context.Context,
	email string,
) (res authdomain.PasswordAuthentication, err error) {
	var emailDoc couchdb.DocumentWithEvents[accountEmailDoc]
	var pwDoc accountPasswordDoc
	_, err1 := r.Connection.Get(r.addrDocId(email), &emailDoc)
	_, err2 := r.Connection.Get(passwordDocId(emailDoc.Document.AccountID), &pwDoc)
	acc, err3 := r.Get(ctx, emailDoc.Document.AccountID)
	if err = errors.Join(err1, err2, err3); err != nil {
		return
	}
	res.Account = acc
	res.PasswordHash = password.HashFromBytes(pwDoc.PasswordHash)
	return
}
