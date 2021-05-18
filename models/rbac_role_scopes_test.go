// Code generated by SQLBoiler 3.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testRbacRoleScopes(t *testing.T) {
	t.Parallel()

	query := RbacRoleScopes()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testRbacRoleScopesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRbacRoleScopesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := RbacRoleScopes().DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRbacRoleScopesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := RbacRoleScopeSlice{o}

	if rowsAff, err := slice.DeleteAll(tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRbacRoleScopesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := RbacRoleScopeExists(tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if RbacRoleScope exists: %s", err)
	}
	if !e {
		t.Errorf("Expected RbacRoleScopeExists to return true, but got false.")
	}
}

func testRbacRoleScopesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	rbacRoleScopeFound, err := FindRbacRoleScope(tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if rbacRoleScopeFound == nil {
		t.Error("want a record, got nil")
	}
}

func testRbacRoleScopesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = RbacRoleScopes().Bind(nil, tx, o); err != nil {
		t.Error(err)
	}
}

func testRbacRoleScopesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := RbacRoleScopes().One(tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testRbacRoleScopesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	rbacRoleScopeOne := &RbacRoleScope{}
	rbacRoleScopeTwo := &RbacRoleScope{}
	if err = randomize.Struct(seed, rbacRoleScopeOne, rbacRoleScopeDBTypes, false, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}
	if err = randomize.Struct(seed, rbacRoleScopeTwo, rbacRoleScopeDBTypes, false, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = rbacRoleScopeOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = rbacRoleScopeTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := RbacRoleScopes().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testRbacRoleScopesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	rbacRoleScopeOne := &RbacRoleScope{}
	rbacRoleScopeTwo := &RbacRoleScope{}
	if err = randomize.Struct(seed, rbacRoleScopeOne, rbacRoleScopeDBTypes, false, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}
	if err = randomize.Struct(seed, rbacRoleScopeTwo, rbacRoleScopeDBTypes, false, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = rbacRoleScopeOne.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = rbacRoleScopeTwo.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func rbacRoleScopeBeforeInsertHook(e boil.Executor, o *RbacRoleScope) error {
	*o = RbacRoleScope{}
	return nil
}

func rbacRoleScopeAfterInsertHook(e boil.Executor, o *RbacRoleScope) error {
	*o = RbacRoleScope{}
	return nil
}

func rbacRoleScopeAfterSelectHook(e boil.Executor, o *RbacRoleScope) error {
	*o = RbacRoleScope{}
	return nil
}

func rbacRoleScopeBeforeUpdateHook(e boil.Executor, o *RbacRoleScope) error {
	*o = RbacRoleScope{}
	return nil
}

func rbacRoleScopeAfterUpdateHook(e boil.Executor, o *RbacRoleScope) error {
	*o = RbacRoleScope{}
	return nil
}

func rbacRoleScopeBeforeDeleteHook(e boil.Executor, o *RbacRoleScope) error {
	*o = RbacRoleScope{}
	return nil
}

func rbacRoleScopeAfterDeleteHook(e boil.Executor, o *RbacRoleScope) error {
	*o = RbacRoleScope{}
	return nil
}

func rbacRoleScopeBeforeUpsertHook(e boil.Executor, o *RbacRoleScope) error {
	*o = RbacRoleScope{}
	return nil
}

func rbacRoleScopeAfterUpsertHook(e boil.Executor, o *RbacRoleScope) error {
	*o = RbacRoleScope{}
	return nil
}

func testRbacRoleScopesHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &RbacRoleScope{}
	o := &RbacRoleScope{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, false); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope object: %s", err)
	}

	AddRbacRoleScopeHook(boil.BeforeInsertHook, rbacRoleScopeBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	rbacRoleScopeBeforeInsertHooks = []RbacRoleScopeHook{}

	AddRbacRoleScopeHook(boil.AfterInsertHook, rbacRoleScopeAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	rbacRoleScopeAfterInsertHooks = []RbacRoleScopeHook{}

	AddRbacRoleScopeHook(boil.AfterSelectHook, rbacRoleScopeAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	rbacRoleScopeAfterSelectHooks = []RbacRoleScopeHook{}

	AddRbacRoleScopeHook(boil.BeforeUpdateHook, rbacRoleScopeBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	rbacRoleScopeBeforeUpdateHooks = []RbacRoleScopeHook{}

	AddRbacRoleScopeHook(boil.AfterUpdateHook, rbacRoleScopeAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	rbacRoleScopeAfterUpdateHooks = []RbacRoleScopeHook{}

	AddRbacRoleScopeHook(boil.BeforeDeleteHook, rbacRoleScopeBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	rbacRoleScopeBeforeDeleteHooks = []RbacRoleScopeHook{}

	AddRbacRoleScopeHook(boil.AfterDeleteHook, rbacRoleScopeAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	rbacRoleScopeAfterDeleteHooks = []RbacRoleScopeHook{}

	AddRbacRoleScopeHook(boil.BeforeUpsertHook, rbacRoleScopeBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	rbacRoleScopeBeforeUpsertHooks = []RbacRoleScopeHook{}

	AddRbacRoleScopeHook(boil.AfterUpsertHook, rbacRoleScopeAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	rbacRoleScopeAfterUpsertHooks = []RbacRoleScopeHook{}
}

func testRbacRoleScopesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRbacRoleScopesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Whitelist(rbacRoleScopeColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRbacRoleScopesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testRbacRoleScopesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := RbacRoleScopeSlice{o}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}

func testRbacRoleScopesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := RbacRoleScopes().All(tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	rbacRoleScopeDBTypes = map[string]string{`ID`: `int`, `CreatedAt`: `timestamp`, `UpdatedAt`: `timestamp`, `DeletedAt`: `timestamp`, `RoleID`: `int`, `ScopeID`: `int`}
	_                    = bytes.MinRead
)

func testRbacRoleScopesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(rbacRoleScopePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(rbacRoleScopeAllColumns) == len(rbacRoleScopePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	if rowsAff, err := o.Update(tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testRbacRoleScopesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(rbacRoleScopeAllColumns) == len(rbacRoleScopePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &RbacRoleScope{}
	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, rbacRoleScopeDBTypes, true, rbacRoleScopePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(rbacRoleScopeAllColumns, rbacRoleScopePrimaryKeyColumns) {
		fields = rbacRoleScopeAllColumns
	} else {
		fields = strmangle.SetComplement(
			rbacRoleScopeAllColumns,
			rbacRoleScopePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := RbacRoleScopeSlice{o}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testRbacRoleScopesUpsert(t *testing.T) {
	t.Parallel()

	if len(rbacRoleScopeAllColumns) == len(rbacRoleScopePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLRbacRoleScopeUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := RbacRoleScope{}
	if err = randomize.Struct(seed, &o, rbacRoleScopeDBTypes, false); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert RbacRoleScope: %s", err)
	}

	count, err := RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, rbacRoleScopeDBTypes, false, rbacRoleScopePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RbacRoleScope struct: %s", err)
	}

	if err = o.Upsert(tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert RbacRoleScope: %s", err)
	}

	count, err = RbacRoleScopes().Count(tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
