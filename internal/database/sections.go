// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Section is an object representing the database table.
type Section struct {
	ID        uint64    `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID    uint64    `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	ListID    uint64    `boil:"list_id" json:"list_id" toml:"list_id" yaml:"list_id"`
	Position  int       `boil:"position" json:"position" toml:"position" yaml:"position"`
	Name      string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *sectionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sectionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SectionColumns = struct {
	ID        string
	UserID    string
	ListID    string
	Position  string
	Name      string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	UserID:    "user_id",
	ListID:    "list_id",
	Position:  "position",
	Name:      "name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var SectionTableColumns = struct {
	ID        string
	UserID    string
	ListID    string
	Position  string
	Name      string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "sections.id",
	UserID:    "sections.user_id",
	ListID:    "sections.list_id",
	Position:  "sections.position",
	Name:      "sections.name",
	CreatedAt: "sections.created_at",
	UpdatedAt: "sections.updated_at",
}

// Generated where

var SectionWhere = struct {
	ID        whereHelperuint64
	UserID    whereHelperuint64
	ListID    whereHelperuint64
	Position  whereHelperint
	Name      whereHelperstring
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	ID:        whereHelperuint64{field: "`sections`.`id`"},
	UserID:    whereHelperuint64{field: "`sections`.`user_id`"},
	ListID:    whereHelperuint64{field: "`sections`.`list_id`"},
	Position:  whereHelperint{field: "`sections`.`position`"},
	Name:      whereHelperstring{field: "`sections`.`name`"},
	CreatedAt: whereHelpertime_Time{field: "`sections`.`created_at`"},
	UpdatedAt: whereHelpertime_Time{field: "`sections`.`updated_at`"},
}

// SectionRels is where relationship names are stored.
var SectionRels = struct {
	List  string
	Items string
}{
	List:  "List",
	Items: "Items",
}

// sectionR is where relationships are stored.
type sectionR struct {
	List  *List     `boil:"List" json:"List" toml:"List" yaml:"List"`
	Items ItemSlice `boil:"Items" json:"Items" toml:"Items" yaml:"Items"`
}

// NewStruct creates a new relationship struct
func (*sectionR) NewStruct() *sectionR {
	return &sectionR{}
}

func (r *sectionR) GetList() *List {
	if r == nil {
		return nil
	}
	return r.List
}

func (r *sectionR) GetItems() ItemSlice {
	if r == nil {
		return nil
	}
	return r.Items
}

// sectionL is where Load methods for each relationship are stored.
type sectionL struct{}

var (
	sectionAllColumns            = []string{"id", "user_id", "list_id", "position", "name", "created_at", "updated_at"}
	sectionColumnsWithoutDefault = []string{"id", "user_id", "list_id", "position", "name"}
	sectionColumnsWithDefault    = []string{"created_at", "updated_at"}
	sectionPrimaryKeyColumns     = []string{"id"}
	sectionGeneratedColumns      = []string{}
)

type (
	// SectionSlice is an alias for a slice of pointers to Section.
	// This should almost always be used instead of []Section.
	SectionSlice []*Section
	// SectionHook is the signature for custom Section hook methods
	SectionHook func(context.Context, boil.ContextExecutor, *Section) error

	sectionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sectionType                 = reflect.TypeOf(&Section{})
	sectionMapping              = queries.MakeStructMapping(sectionType)
	sectionPrimaryKeyMapping, _ = queries.BindMapping(sectionType, sectionMapping, sectionPrimaryKeyColumns)
	sectionInsertCacheMut       sync.RWMutex
	sectionInsertCache          = make(map[string]insertCache)
	sectionUpdateCacheMut       sync.RWMutex
	sectionUpdateCache          = make(map[string]updateCache)
	sectionUpsertCacheMut       sync.RWMutex
	sectionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var sectionAfterSelectHooks []SectionHook

var sectionBeforeInsertHooks []SectionHook
var sectionAfterInsertHooks []SectionHook

var sectionBeforeUpdateHooks []SectionHook
var sectionAfterUpdateHooks []SectionHook

var sectionBeforeDeleteHooks []SectionHook
var sectionAfterDeleteHooks []SectionHook

var sectionBeforeUpsertHooks []SectionHook
var sectionAfterUpsertHooks []SectionHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Section) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sectionAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Section) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sectionBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Section) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sectionAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Section) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sectionBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Section) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sectionAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Section) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sectionBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Section) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sectionAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Section) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sectionBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Section) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range sectionAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSectionHook registers your hook function for all future operations.
func AddSectionHook(hookPoint boil.HookPoint, sectionHook SectionHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		sectionAfterSelectHooks = append(sectionAfterSelectHooks, sectionHook)
	case boil.BeforeInsertHook:
		sectionBeforeInsertHooks = append(sectionBeforeInsertHooks, sectionHook)
	case boil.AfterInsertHook:
		sectionAfterInsertHooks = append(sectionAfterInsertHooks, sectionHook)
	case boil.BeforeUpdateHook:
		sectionBeforeUpdateHooks = append(sectionBeforeUpdateHooks, sectionHook)
	case boil.AfterUpdateHook:
		sectionAfterUpdateHooks = append(sectionAfterUpdateHooks, sectionHook)
	case boil.BeforeDeleteHook:
		sectionBeforeDeleteHooks = append(sectionBeforeDeleteHooks, sectionHook)
	case boil.AfterDeleteHook:
		sectionAfterDeleteHooks = append(sectionAfterDeleteHooks, sectionHook)
	case boil.BeforeUpsertHook:
		sectionBeforeUpsertHooks = append(sectionBeforeUpsertHooks, sectionHook)
	case boil.AfterUpsertHook:
		sectionAfterUpsertHooks = append(sectionAfterUpsertHooks, sectionHook)
	}
}

// One returns a single section record from the query.
func (q sectionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Section, error) {
	o := &Section{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: failed to execute a one query for sections")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Section records from the query.
func (q sectionQuery) All(ctx context.Context, exec boil.ContextExecutor) (SectionSlice, error) {
	var o []*Section

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "database: failed to assign all query results to Section slice")
	}

	if len(sectionAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Section records in the query.
func (q sectionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to count sections rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q sectionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "database: failed to check if sections exists")
	}

	return count > 0, nil
}

// List pointed to by the foreign key.
func (o *Section) List(mods ...qm.QueryMod) listQuery {
	queryMods := []qm.QueryMod{
		qm.Where("`id` = ?", o.ListID),
	}

	queryMods = append(queryMods, mods...)

	return Lists(queryMods...)
}

// Items retrieves all the item's Items with an executor.
func (o *Section) Items(mods ...qm.QueryMod) itemQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("`items`.`section_id`=?", o.ID),
	)

	return Items(queryMods...)
}

// LoadList allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (sectionL) LoadList(ctx context.Context, e boil.ContextExecutor, singular bool, maybeSection interface{}, mods queries.Applicator) error {
	var slice []*Section
	var object *Section

	if singular {
		var ok bool
		object, ok = maybeSection.(*Section)
		if !ok {
			object = new(Section)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeSection)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeSection))
			}
		}
	} else {
		s, ok := maybeSection.(*[]*Section)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeSection)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeSection))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &sectionR{}
		}
		args = append(args, object.ListID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &sectionR{}
			}

			for _, a := range args {
				if a == obj.ListID {
					continue Outer
				}
			}

			args = append(args, obj.ListID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`lists`),
		qm.WhereIn(`lists.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load List")
	}

	var resultSlice []*List
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice List")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for lists")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for lists")
	}

	if len(listAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.List = foreign
		if foreign.R == nil {
			foreign.R = &listR{}
		}
		foreign.R.Sections = append(foreign.R.Sections, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ListID == foreign.ID {
				local.R.List = foreign
				if foreign.R == nil {
					foreign.R = &listR{}
				}
				foreign.R.Sections = append(foreign.R.Sections, local)
				break
			}
		}
	}

	return nil
}

// LoadItems allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (sectionL) LoadItems(ctx context.Context, e boil.ContextExecutor, singular bool, maybeSection interface{}, mods queries.Applicator) error {
	var slice []*Section
	var object *Section

	if singular {
		var ok bool
		object, ok = maybeSection.(*Section)
		if !ok {
			object = new(Section)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeSection)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeSection))
			}
		}
	} else {
		s, ok := maybeSection.(*[]*Section)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeSection)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeSection))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &sectionR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &sectionR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`items`),
		qm.WhereIn(`items.section_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load items")
	}

	var resultSlice []*Item
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice items")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on items")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for items")
	}

	if len(itemAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Items = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &itemR{}
			}
			foreign.R.Section = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.SectionID) {
				local.R.Items = append(local.R.Items, foreign)
				if foreign.R == nil {
					foreign.R = &itemR{}
				}
				foreign.R.Section = local
				break
			}
		}
	}

	return nil
}

// SetList of the section to the related item.
// Sets o.R.List to related.
// Adds o to related.R.Sections.
func (o *Section) SetList(ctx context.Context, exec boil.ContextExecutor, insert bool, related *List) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `sections` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"list_id"}),
		strmangle.WhereClause("`", "`", 0, sectionPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ListID = related.ID
	if o.R == nil {
		o.R = &sectionR{
			List: related,
		}
	} else {
		o.R.List = related
	}

	if related.R == nil {
		related.R = &listR{
			Sections: SectionSlice{o},
		}
	} else {
		related.R.Sections = append(related.R.Sections, o)
	}

	return nil
}

// AddItems adds the given related objects to the existing relationships
// of the section, optionally inserting them as new records.
// Appends related to o.R.Items.
// Sets related.R.Section appropriately.
func (o *Section) AddItems(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Item) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.SectionID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE `items` SET %s WHERE %s",
				strmangle.SetParamNames("`", "`", 0, []string{"section_id"}),
				strmangle.WhereClause("`", "`", 0, itemPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.SectionID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &sectionR{
			Items: related,
		}
	} else {
		o.R.Items = append(o.R.Items, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &itemR{
				Section: o,
			}
		} else {
			rel.R.Section = o
		}
	}
	return nil
}

// SetItems removes all previously related items of the
// section replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Section's Items accordingly.
// Replaces o.R.Items with related.
// Sets related.R.Section's Items accordingly.
func (o *Section) SetItems(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Item) error {
	query := "update `items` set `section_id` = null where `section_id` = ?"
	values := []interface{}{o.ID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.Items {
			queries.SetScanner(&rel.SectionID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Section = nil
		}
		o.R.Items = nil
	}

	return o.AddItems(ctx, exec, insert, related...)
}

// RemoveItems relationships from objects passed in.
// Removes related items from R.Items (uses pointer comparison, removal does not keep order)
// Sets related.R.Section.
func (o *Section) RemoveItems(ctx context.Context, exec boil.ContextExecutor, related ...*Item) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.SectionID, nil)
		if rel.R != nil {
			rel.R.Section = nil
		}
		if _, err = rel.Update(ctx, exec, boil.Whitelist("section_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Items {
			if rel != ri {
				continue
			}

			ln := len(o.R.Items)
			if ln > 1 && i < ln-1 {
				o.R.Items[i] = o.R.Items[ln-1]
			}
			o.R.Items = o.R.Items[:ln-1]
			break
		}
	}

	return nil
}

// Sections retrieves all the records using an executor.
func Sections(mods ...qm.QueryMod) sectionQuery {
	mods = append(mods, qm.From("`sections`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`sections`.*"})
	}

	return sectionQuery{q}
}

// FindSection retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSection(ctx context.Context, exec boil.ContextExecutor, iD uint64, selectCols ...string) (*Section, error) {
	sectionObj := &Section{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `sections` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, sectionObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "database: unable to select from sections")
	}

	if err = sectionObj.doAfterSelectHooks(ctx, exec); err != nil {
		return sectionObj, err
	}

	return sectionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Section) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("database: no sections provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(sectionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	sectionInsertCacheMut.RLock()
	cache, cached := sectionInsertCache[key]
	sectionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			sectionAllColumns,
			sectionColumnsWithDefault,
			sectionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(sectionType, sectionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sectionType, sectionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `sections` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `sections` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `sections` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, sectionPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "database: unable to insert into sections")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "database: unable to populate default values for sections")
	}

CacheNoHooks:
	if !cached {
		sectionInsertCacheMut.Lock()
		sectionInsertCache[key] = cache
		sectionInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Section.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Section) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	sectionUpdateCacheMut.RLock()
	cache, cached := sectionUpdateCache[key]
	sectionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			sectionAllColumns,
			sectionPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("database: unable to update sections, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `sections` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, sectionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sectionType, sectionMapping, append(wl, sectionPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update sections row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by update for sections")
	}

	if !cached {
		sectionUpdateCacheMut.Lock()
		sectionUpdateCache[key] = cache
		sectionUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q sectionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all for sections")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected for sections")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SectionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("database: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `sections` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, sectionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to update all in section slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to retrieve rows affected all in update all section")
	}
	return rowsAff, nil
}

var mySQLSectionUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Section) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("database: no sections provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(sectionColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLSectionUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	sectionUpsertCacheMut.RLock()
	cache, cached := sectionUpsertCache[key]
	sectionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			sectionAllColumns,
			sectionColumnsWithDefault,
			sectionColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			sectionAllColumns,
			sectionPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("database: unable to upsert sections, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`sections`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `sections` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(sectionType, sectionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sectionType, sectionMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "database: unable to upsert for sections")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(sectionType, sectionMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "database: unable to retrieve unique values for sections")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "database: unable to populate default values for sections")
	}

CacheNoHooks:
	if !cached {
		sectionUpsertCacheMut.Lock()
		sectionUpsertCache[key] = cache
		sectionUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Section record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Section) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("database: no Section provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sectionPrimaryKeyMapping)
	sql := "DELETE FROM `sections` WHERE `id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete from sections")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by delete for sections")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q sectionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("database: no sectionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from sections")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for sections")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SectionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(sectionBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `sections` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, sectionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "database: unable to delete all from section slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "database: failed to get rows affected by deleteall for sections")
	}

	if len(sectionAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Section) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSection(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SectionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SectionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `sections`.* FROM `sections` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, sectionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "database: unable to reload all in SectionSlice")
	}

	*o = slice

	return nil
}

// SectionExists checks if the Section row exists.
func SectionExists(ctx context.Context, exec boil.ContextExecutor, iD uint64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `sections` where `id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "database: unable to check if sections exists")
	}

	return exists, nil
}

// Exists checks if the Section row exists.
func (o *Section) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return SectionExists(ctx, exec, o.ID)
}
