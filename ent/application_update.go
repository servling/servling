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
	"github.com/servling/servling/ent/service"
	"github.com/servling/servling/ent/template"
)

// ApplicationUpdate is the builder for updating Application entities.
type ApplicationUpdate struct {
	config
	hooks    []Hook
	mutation *ApplicationMutation
}

// Where appends a list predicates to the ApplicationUpdate builder.
func (au *ApplicationUpdate) Where(ps ...predicate.Application) *ApplicationUpdate {
	au.mutation.Where(ps...)
	return au
}

// SetName sets the "name" field.
func (au *ApplicationUpdate) SetName(s string) *ApplicationUpdate {
	au.mutation.SetName(s)
	return au
}

// SetNillableName sets the "name" field if the given value is not nil.
func (au *ApplicationUpdate) SetNillableName(s *string) *ApplicationUpdate {
	if s != nil {
		au.SetName(*s)
	}
	return au
}

// SetDescription sets the "description" field.
func (au *ApplicationUpdate) SetDescription(s string) *ApplicationUpdate {
	au.mutation.SetDescription(s)
	return au
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (au *ApplicationUpdate) SetNillableDescription(s *string) *ApplicationUpdate {
	if s != nil {
		au.SetDescription(*s)
	}
	return au
}

// SetImageURL sets the "image_url" field.
func (au *ApplicationUpdate) SetImageURL(s string) *ApplicationUpdate {
	au.mutation.SetImageURL(s)
	return au
}

// SetNillableImageURL sets the "image_url" field if the given value is not nil.
func (au *ApplicationUpdate) SetNillableImageURL(s *string) *ApplicationUpdate {
	if s != nil {
		au.SetImageURL(*s)
	}
	return au
}

// ClearImageURL clears the value of the "image_url" field.
func (au *ApplicationUpdate) ClearImageURL() *ApplicationUpdate {
	au.mutation.ClearImageURL()
	return au
}

// SetStatus sets the "status" field.
func (au *ApplicationUpdate) SetStatus(s string) *ApplicationUpdate {
	au.mutation.SetStatus(s)
	return au
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (au *ApplicationUpdate) SetNillableStatus(s *string) *ApplicationUpdate {
	if s != nil {
		au.SetStatus(*s)
	}
	return au
}

// SetError sets the "error" field.
func (au *ApplicationUpdate) SetError(s string) *ApplicationUpdate {
	au.mutation.SetError(s)
	return au
}

// SetNillableError sets the "error" field if the given value is not nil.
func (au *ApplicationUpdate) SetNillableError(s *string) *ApplicationUpdate {
	if s != nil {
		au.SetError(*s)
	}
	return au
}

// ClearError clears the value of the "error" field.
func (au *ApplicationUpdate) ClearError() *ApplicationUpdate {
	au.mutation.ClearError()
	return au
}

// SetUpdatedAt sets the "updated_at" field.
func (au *ApplicationUpdate) SetUpdatedAt(t time.Time) *ApplicationUpdate {
	au.mutation.SetUpdatedAt(t)
	return au
}

// AddServiceIDs adds the "services" edge to the Service entity by IDs.
func (au *ApplicationUpdate) AddServiceIDs(ids ...string) *ApplicationUpdate {
	au.mutation.AddServiceIDs(ids...)
	return au
}

// AddServices adds the "services" edges to the Service entity.
func (au *ApplicationUpdate) AddServices(s ...*Service) *ApplicationUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return au.AddServiceIDs(ids...)
}

// SetTemplateID sets the "template" edge to the Template entity by ID.
func (au *ApplicationUpdate) SetTemplateID(id string) *ApplicationUpdate {
	au.mutation.SetTemplateID(id)
	return au
}

// SetNillableTemplateID sets the "template" edge to the Template entity by ID if the given value is not nil.
func (au *ApplicationUpdate) SetNillableTemplateID(id *string) *ApplicationUpdate {
	if id != nil {
		au = au.SetTemplateID(*id)
	}
	return au
}

// SetTemplate sets the "template" edge to the Template entity.
func (au *ApplicationUpdate) SetTemplate(t *Template) *ApplicationUpdate {
	return au.SetTemplateID(t.ID)
}

// Mutation returns the ApplicationMutation object of the builder.
func (au *ApplicationUpdate) Mutation() *ApplicationMutation {
	return au.mutation
}

// ClearServices clears all "services" edges to the Service entity.
func (au *ApplicationUpdate) ClearServices() *ApplicationUpdate {
	au.mutation.ClearServices()
	return au
}

// RemoveServiceIDs removes the "services" edge to Service entities by IDs.
func (au *ApplicationUpdate) RemoveServiceIDs(ids ...string) *ApplicationUpdate {
	au.mutation.RemoveServiceIDs(ids...)
	return au
}

// RemoveServices removes "services" edges to Service entities.
func (au *ApplicationUpdate) RemoveServices(s ...*Service) *ApplicationUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return au.RemoveServiceIDs(ids...)
}

// ClearTemplate clears the "template" edge to the Template entity.
func (au *ApplicationUpdate) ClearTemplate() *ApplicationUpdate {
	au.mutation.ClearTemplate()
	return au
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (au *ApplicationUpdate) Save(ctx context.Context) (int, error) {
	au.defaults()
	return withHooks(ctx, au.sqlSave, au.mutation, au.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (au *ApplicationUpdate) SaveX(ctx context.Context) int {
	affected, err := au.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (au *ApplicationUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (au *ApplicationUpdate) ExecX(ctx context.Context) {
	if err := au.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (au *ApplicationUpdate) defaults() {
	if _, ok := au.mutation.UpdatedAt(); !ok {
		v := application.UpdateDefaultUpdatedAt()
		au.mutation.SetUpdatedAt(v)
	}
}

func (au *ApplicationUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(application.Table, application.Columns, sqlgraph.NewFieldSpec(application.FieldID, field.TypeString))
	if ps := au.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := au.mutation.Name(); ok {
		_spec.SetField(application.FieldName, field.TypeString, value)
	}
	if value, ok := au.mutation.Description(); ok {
		_spec.SetField(application.FieldDescription, field.TypeString, value)
	}
	if value, ok := au.mutation.ImageURL(); ok {
		_spec.SetField(application.FieldImageURL, field.TypeString, value)
	}
	if au.mutation.ImageURLCleared() {
		_spec.ClearField(application.FieldImageURL, field.TypeString)
	}
	if value, ok := au.mutation.Status(); ok {
		_spec.SetField(application.FieldStatus, field.TypeString, value)
	}
	if value, ok := au.mutation.Error(); ok {
		_spec.SetField(application.FieldError, field.TypeString, value)
	}
	if au.mutation.ErrorCleared() {
		_spec.ClearField(application.FieldError, field.TypeString)
	}
	if value, ok := au.mutation.UpdatedAt(); ok {
		_spec.SetField(application.FieldUpdatedAt, field.TypeTime, value)
	}
	if au.mutation.ServicesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.ServicesTable,
			Columns: []string{application.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.RemovedServicesIDs(); len(nodes) > 0 && !au.mutation.ServicesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.ServicesTable,
			Columns: []string{application.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.ServicesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.ServicesTable,
			Columns: []string{application.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if au.mutation.TemplateCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   application.TemplateTable,
			Columns: []string{application.TemplateColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.TemplateIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   application.TemplateTable,
			Columns: []string{application.TemplateColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, au.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{application.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	au.mutation.done = true
	return n, nil
}

// ApplicationUpdateOne is the builder for updating a single Application entity.
type ApplicationUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ApplicationMutation
}

// SetName sets the "name" field.
func (auo *ApplicationUpdateOne) SetName(s string) *ApplicationUpdateOne {
	auo.mutation.SetName(s)
	return auo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (auo *ApplicationUpdateOne) SetNillableName(s *string) *ApplicationUpdateOne {
	if s != nil {
		auo.SetName(*s)
	}
	return auo
}

// SetDescription sets the "description" field.
func (auo *ApplicationUpdateOne) SetDescription(s string) *ApplicationUpdateOne {
	auo.mutation.SetDescription(s)
	return auo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (auo *ApplicationUpdateOne) SetNillableDescription(s *string) *ApplicationUpdateOne {
	if s != nil {
		auo.SetDescription(*s)
	}
	return auo
}

// SetImageURL sets the "image_url" field.
func (auo *ApplicationUpdateOne) SetImageURL(s string) *ApplicationUpdateOne {
	auo.mutation.SetImageURL(s)
	return auo
}

// SetNillableImageURL sets the "image_url" field if the given value is not nil.
func (auo *ApplicationUpdateOne) SetNillableImageURL(s *string) *ApplicationUpdateOne {
	if s != nil {
		auo.SetImageURL(*s)
	}
	return auo
}

// ClearImageURL clears the value of the "image_url" field.
func (auo *ApplicationUpdateOne) ClearImageURL() *ApplicationUpdateOne {
	auo.mutation.ClearImageURL()
	return auo
}

// SetStatus sets the "status" field.
func (auo *ApplicationUpdateOne) SetStatus(s string) *ApplicationUpdateOne {
	auo.mutation.SetStatus(s)
	return auo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (auo *ApplicationUpdateOne) SetNillableStatus(s *string) *ApplicationUpdateOne {
	if s != nil {
		auo.SetStatus(*s)
	}
	return auo
}

// SetError sets the "error" field.
func (auo *ApplicationUpdateOne) SetError(s string) *ApplicationUpdateOne {
	auo.mutation.SetError(s)
	return auo
}

// SetNillableError sets the "error" field if the given value is not nil.
func (auo *ApplicationUpdateOne) SetNillableError(s *string) *ApplicationUpdateOne {
	if s != nil {
		auo.SetError(*s)
	}
	return auo
}

// ClearError clears the value of the "error" field.
func (auo *ApplicationUpdateOne) ClearError() *ApplicationUpdateOne {
	auo.mutation.ClearError()
	return auo
}

// SetUpdatedAt sets the "updated_at" field.
func (auo *ApplicationUpdateOne) SetUpdatedAt(t time.Time) *ApplicationUpdateOne {
	auo.mutation.SetUpdatedAt(t)
	return auo
}

// AddServiceIDs adds the "services" edge to the Service entity by IDs.
func (auo *ApplicationUpdateOne) AddServiceIDs(ids ...string) *ApplicationUpdateOne {
	auo.mutation.AddServiceIDs(ids...)
	return auo
}

// AddServices adds the "services" edges to the Service entity.
func (auo *ApplicationUpdateOne) AddServices(s ...*Service) *ApplicationUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return auo.AddServiceIDs(ids...)
}

// SetTemplateID sets the "template" edge to the Template entity by ID.
func (auo *ApplicationUpdateOne) SetTemplateID(id string) *ApplicationUpdateOne {
	auo.mutation.SetTemplateID(id)
	return auo
}

// SetNillableTemplateID sets the "template" edge to the Template entity by ID if the given value is not nil.
func (auo *ApplicationUpdateOne) SetNillableTemplateID(id *string) *ApplicationUpdateOne {
	if id != nil {
		auo = auo.SetTemplateID(*id)
	}
	return auo
}

// SetTemplate sets the "template" edge to the Template entity.
func (auo *ApplicationUpdateOne) SetTemplate(t *Template) *ApplicationUpdateOne {
	return auo.SetTemplateID(t.ID)
}

// Mutation returns the ApplicationMutation object of the builder.
func (auo *ApplicationUpdateOne) Mutation() *ApplicationMutation {
	return auo.mutation
}

// ClearServices clears all "services" edges to the Service entity.
func (auo *ApplicationUpdateOne) ClearServices() *ApplicationUpdateOne {
	auo.mutation.ClearServices()
	return auo
}

// RemoveServiceIDs removes the "services" edge to Service entities by IDs.
func (auo *ApplicationUpdateOne) RemoveServiceIDs(ids ...string) *ApplicationUpdateOne {
	auo.mutation.RemoveServiceIDs(ids...)
	return auo
}

// RemoveServices removes "services" edges to Service entities.
func (auo *ApplicationUpdateOne) RemoveServices(s ...*Service) *ApplicationUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return auo.RemoveServiceIDs(ids...)
}

// ClearTemplate clears the "template" edge to the Template entity.
func (auo *ApplicationUpdateOne) ClearTemplate() *ApplicationUpdateOne {
	auo.mutation.ClearTemplate()
	return auo
}

// Where appends a list predicates to the ApplicationUpdate builder.
func (auo *ApplicationUpdateOne) Where(ps ...predicate.Application) *ApplicationUpdateOne {
	auo.mutation.Where(ps...)
	return auo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (auo *ApplicationUpdateOne) Select(field string, fields ...string) *ApplicationUpdateOne {
	auo.fields = append([]string{field}, fields...)
	return auo
}

// Save executes the query and returns the updated Application entity.
func (auo *ApplicationUpdateOne) Save(ctx context.Context) (*Application, error) {
	auo.defaults()
	return withHooks(ctx, auo.sqlSave, auo.mutation, auo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (auo *ApplicationUpdateOne) SaveX(ctx context.Context) *Application {
	node, err := auo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (auo *ApplicationUpdateOne) Exec(ctx context.Context) error {
	_, err := auo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auo *ApplicationUpdateOne) ExecX(ctx context.Context) {
	if err := auo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (auo *ApplicationUpdateOne) defaults() {
	if _, ok := auo.mutation.UpdatedAt(); !ok {
		v := application.UpdateDefaultUpdatedAt()
		auo.mutation.SetUpdatedAt(v)
	}
}

func (auo *ApplicationUpdateOne) sqlSave(ctx context.Context) (_node *Application, err error) {
	_spec := sqlgraph.NewUpdateSpec(application.Table, application.Columns, sqlgraph.NewFieldSpec(application.FieldID, field.TypeString))
	id, ok := auo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Application.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := auo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, application.FieldID)
		for _, f := range fields {
			if !application.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != application.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := auo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := auo.mutation.Name(); ok {
		_spec.SetField(application.FieldName, field.TypeString, value)
	}
	if value, ok := auo.mutation.Description(); ok {
		_spec.SetField(application.FieldDescription, field.TypeString, value)
	}
	if value, ok := auo.mutation.ImageURL(); ok {
		_spec.SetField(application.FieldImageURL, field.TypeString, value)
	}
	if auo.mutation.ImageURLCleared() {
		_spec.ClearField(application.FieldImageURL, field.TypeString)
	}
	if value, ok := auo.mutation.Status(); ok {
		_spec.SetField(application.FieldStatus, field.TypeString, value)
	}
	if value, ok := auo.mutation.Error(); ok {
		_spec.SetField(application.FieldError, field.TypeString, value)
	}
	if auo.mutation.ErrorCleared() {
		_spec.ClearField(application.FieldError, field.TypeString)
	}
	if value, ok := auo.mutation.UpdatedAt(); ok {
		_spec.SetField(application.FieldUpdatedAt, field.TypeTime, value)
	}
	if auo.mutation.ServicesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.ServicesTable,
			Columns: []string{application.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.RemovedServicesIDs(); len(nodes) > 0 && !auo.mutation.ServicesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.ServicesTable,
			Columns: []string{application.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.ServicesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.ServicesTable,
			Columns: []string{application.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if auo.mutation.TemplateCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   application.TemplateTable,
			Columns: []string{application.TemplateColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.TemplateIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   application.TemplateTable,
			Columns: []string{application.TemplateColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Application{config: auo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, auo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{application.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	auo.mutation.done = true
	return _node, nil
}
