// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/servling/servling/ent/application"
	"github.com/servling/servling/ent/predicate"
	"github.com/servling/servling/ent/template"
)

// TemplateUpdate is the builder for updating Template entities.
type TemplateUpdate struct {
	config
	hooks    []Hook
	mutation *TemplateMutation
}

// Where appends a list predicates to the TemplateUpdate builder.
func (tu *TemplateUpdate) Where(ps ...predicate.Template) *TemplateUpdate {
	tu.mutation.Where(ps...)
	return tu
}

// SetName sets the "name" field.
func (tu *TemplateUpdate) SetName(s string) *TemplateUpdate {
	tu.mutation.SetName(s)
	return tu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tu *TemplateUpdate) SetNillableName(s *string) *TemplateUpdate {
	if s != nil {
		tu.SetName(*s)
	}
	return tu
}

// SetUpdatedAt sets the "updated_at" field.
func (tu *TemplateUpdate) SetUpdatedAt(t time.Time) *TemplateUpdate {
	tu.mutation.SetUpdatedAt(t)
	return tu
}

// AddApplicationIDs adds the "applications" edge to the Application entity by IDs.
func (tu *TemplateUpdate) AddApplicationIDs(ids ...string) *TemplateUpdate {
	tu.mutation.AddApplicationIDs(ids...)
	return tu
}

// AddApplications adds the "applications" edges to the Application entity.
func (tu *TemplateUpdate) AddApplications(a ...*Application) *TemplateUpdate {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return tu.AddApplicationIDs(ids...)
}

// Mutation returns the TemplateMutation object of the builder.
func (tu *TemplateUpdate) Mutation() *TemplateMutation {
	return tu.mutation
}

// ClearApplications clears all "applications" edges to the Application entity.
func (tu *TemplateUpdate) ClearApplications() *TemplateUpdate {
	tu.mutation.ClearApplications()
	return tu
}

// RemoveApplicationIDs removes the "applications" edge to Application entities by IDs.
func (tu *TemplateUpdate) RemoveApplicationIDs(ids ...string) *TemplateUpdate {
	tu.mutation.RemoveApplicationIDs(ids...)
	return tu
}

// RemoveApplications removes "applications" edges to Application entities.
func (tu *TemplateUpdate) RemoveApplications(a ...*Application) *TemplateUpdate {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return tu.RemoveApplicationIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TemplateUpdate) Save(ctx context.Context) (int, error) {
	tu.defaults()
	return withHooks(ctx, tu.sqlSave, tu.mutation, tu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TemplateUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TemplateUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TemplateUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tu *TemplateUpdate) defaults() {
	if _, ok := tu.mutation.UpdatedAt(); !ok {
		v := template.UpdateDefaultUpdatedAt()
		tu.mutation.SetUpdatedAt(v)
	}
}

func (tu *TemplateUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(template.Table, template.Columns, sqlgraph.NewFieldSpec(template.FieldID, field.TypeString))
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.Name(); ok {
		_spec.SetField(template.FieldName, field.TypeString, value)
	}
	if value, ok := tu.mutation.UpdatedAt(); ok {
		_spec.SetField(template.FieldUpdatedAt, field.TypeTime, value)
	}
	if tu.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.ApplicationsTable,
			Columns: []string{template.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(application.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedApplicationsIDs(); len(nodes) > 0 && !tu.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.ApplicationsTable,
			Columns: []string{template.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(application.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.ApplicationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.ApplicationsTable,
			Columns: []string{template.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(application.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{template.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tu.mutation.done = true
	return n, nil
}

// TemplateUpdateOne is the builder for updating a single Template entity.
type TemplateUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TemplateMutation
}

// SetName sets the "name" field.
func (tuo *TemplateUpdateOne) SetName(s string) *TemplateUpdateOne {
	tuo.mutation.SetName(s)
	return tuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tuo *TemplateUpdateOne) SetNillableName(s *string) *TemplateUpdateOne {
	if s != nil {
		tuo.SetName(*s)
	}
	return tuo
}

// SetUpdatedAt sets the "updated_at" field.
func (tuo *TemplateUpdateOne) SetUpdatedAt(t time.Time) *TemplateUpdateOne {
	tuo.mutation.SetUpdatedAt(t)
	return tuo
}

// AddApplicationIDs adds the "applications" edge to the Application entity by IDs.
func (tuo *TemplateUpdateOne) AddApplicationIDs(ids ...string) *TemplateUpdateOne {
	tuo.mutation.AddApplicationIDs(ids...)
	return tuo
}

// AddApplications adds the "applications" edges to the Application entity.
func (tuo *TemplateUpdateOne) AddApplications(a ...*Application) *TemplateUpdateOne {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return tuo.AddApplicationIDs(ids...)
}

// Mutation returns the TemplateMutation object of the builder.
func (tuo *TemplateUpdateOne) Mutation() *TemplateMutation {
	return tuo.mutation
}

// ClearApplications clears all "applications" edges to the Application entity.
func (tuo *TemplateUpdateOne) ClearApplications() *TemplateUpdateOne {
	tuo.mutation.ClearApplications()
	return tuo
}

// RemoveApplicationIDs removes the "applications" edge to Application entities by IDs.
func (tuo *TemplateUpdateOne) RemoveApplicationIDs(ids ...string) *TemplateUpdateOne {
	tuo.mutation.RemoveApplicationIDs(ids...)
	return tuo
}

// RemoveApplications removes "applications" edges to Application entities.
func (tuo *TemplateUpdateOne) RemoveApplications(a ...*Application) *TemplateUpdateOne {
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return tuo.RemoveApplicationIDs(ids...)
}

// Where appends a list predicates to the TemplateUpdate builder.
func (tuo *TemplateUpdateOne) Where(ps ...predicate.Template) *TemplateUpdateOne {
	tuo.mutation.Where(ps...)
	return tuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tuo *TemplateUpdateOne) Select(field string, fields ...string) *TemplateUpdateOne {
	tuo.fields = append([]string{field}, fields...)
	return tuo
}

// Save executes the query and returns the updated Template entity.
func (tuo *TemplateUpdateOne) Save(ctx context.Context) (*Template, error) {
	tuo.defaults()
	return withHooks(ctx, tuo.sqlSave, tuo.mutation, tuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TemplateUpdateOne) SaveX(ctx context.Context) *Template {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TemplateUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TemplateUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tuo *TemplateUpdateOne) defaults() {
	if _, ok := tuo.mutation.UpdatedAt(); !ok {
		v := template.UpdateDefaultUpdatedAt()
		tuo.mutation.SetUpdatedAt(v)
	}
}

func (tuo *TemplateUpdateOne) sqlSave(ctx context.Context) (_node *Template, err error) {
	_spec := sqlgraph.NewUpdateSpec(template.Table, template.Columns, sqlgraph.NewFieldSpec(template.FieldID, field.TypeString))
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Template.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, template.FieldID)
		for _, f := range fields {
			if !template.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != template.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.Name(); ok {
		_spec.SetField(template.FieldName, field.TypeString, value)
	}
	if value, ok := tuo.mutation.UpdatedAt(); ok {
		_spec.SetField(template.FieldUpdatedAt, field.TypeTime, value)
	}
	if tuo.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.ApplicationsTable,
			Columns: []string{template.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(application.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedApplicationsIDs(); len(nodes) > 0 && !tuo.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.ApplicationsTable,
			Columns: []string{template.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(application.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.ApplicationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.ApplicationsTable,
			Columns: []string{template.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(application.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Template{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{template.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tuo.mutation.done = true
	return _node, nil
}
