package api

import (
	"context"
	"fmt"
	"log"

	"github.com/pscn/gqlgen-todos/graph"
	"github.com/pscn/gqlgen-todos/model"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	urls         map[string]model.URL
	vendors      map[string]model.Vendor
	concentrates map[string]model.Concentrate
}

func (r *Resolver) addURL(url *model.URL) error {
	if r.urls == nil {
		r.urls = make(map[string]model.URL, 1)
	}
	r.urls[url.ID] = *url
	return nil
}

func (r *Resolver) getURL(id string) *model.URL {
	if result, ok := r.urls[id]; ok {
		return &result
	}
	return nil
}

func (r *Resolver) getURLs(ids []*string) []*model.URL {
	result := make([]*model.URL, len(ids))
	for idx, id := range ids {
		result[idx] = r.getURL(*id)
	}
	return result
}

func (r *Resolver) addVendor(vendor *model.Vendor) error {
	if r.vendors == nil {
		r.vendors = make(map[string]model.Vendor, 1)
	}
	r.vendors[vendor.ID] = *vendor
	return nil
}

func (r *Resolver) getVendor(id string) *model.Vendor {
	if result, ok := r.vendors[id]; ok {
		return &result
	}
	return nil
}

func (r *Resolver) addConcentrate(concentrate *model.Concentrate) error {
	if r.concentrates == nil {
		r.concentrates = make(map[string]model.Concentrate, 1)
	}
	r.concentrates[concentrate.ID] = *concentrate
	return nil
}

func (r *Resolver) getConcentrate(id string) *model.Concentrate {
	if result, ok := r.concentrates[id]; ok {
		return &result
	}
	return nil
}

func (r *Resolver) Concentrate() graph.ConcentrateResolver {
	return &concentrateResolver{r}
}
func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Vendor() graph.VendorResolver {
	return &vendorResolver{r}
}

type concentrateResolver struct{ *Resolver }

func (r *concentrateResolver) Vendor(ctx context.Context, obj *model.Concentrate) (*model.Vendor, error) {
	return r.getVendor(obj.VendorID), nil
}
func (r *concentrateResolver) Urls(ctx context.Context, obj *model.Concentrate) ([]*model.URL, error) {
	return r.getURLs(obj.URLIDs), nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateVendor(ctx context.Context, input model.NewVendor) (*model.Vendor, error) {
	id := fmt.Sprintf("%s-%s", input.Code, input.Name)
	vendor := &model.Vendor{
		ID:   id,
		Name: input.Name,
		Code: input.Code,
	}
	r.addVendor(vendor)
	return vendor, nil
}
func (r *mutationResolver) AddVendorURL(ctx context.Context, input model.NewVendorURL) (*model.URL, error) {
	vendor := r.getVendor(input.VendorID)
	if vendor == nil {
		return nil, fmt.Errorf("vendor not found")
	}
	id := input.URL
	url := &model.URL{
		ID:          id,
		Description: input.Description,
		URL:         input.URL,
	}
	r.addURL(url)
	vendor.URLIDs = append(vendor.URLIDs, &id)
	return url, nil
}
func (r *mutationResolver) CreateConcentrate(ctx context.Context, input model.NewConcentrate) (*model.Concentrate, error) {
	vendor := r.getVendor(input.VendorID)
	if vendor == nil {
		return nil, fmt.Errorf("vendor not found")
	}
	id := fmt.Sprintf("%s-%s", vendor.Code, input.Name)
	concentrate := &model.Concentrate{
		ID:       id,
		Name:     input.Name,
		VendorID: vendor.ID,
	}
	r.addConcentrate(concentrate)
	return concentrate, nil
}
func (r *mutationResolver) AddConcentratURL(ctx context.Context, input model.NewConcentrateURL) (*model.URL, error) {
	concentrate := r.getConcentrate(input.ConcentrateID)
	if concentrate == nil {
		return nil, fmt.Errorf("concentrate not found")
	}
	id := input.URL
	url := &model.URL{
		ID:          id,
		Description: input.Description,
		URL:         input.URL,
	}
	r.addURL(url)
	concentrate.URLIDs = append(concentrate.URLIDs, &id)
	return url, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Vendors(ctx context.Context) ([]model.Vendor, error) {
	log.Printf("%+v", r.vendors)
	result := make([]model.Vendor, len(r.vendors))
	i := 0
	for _, vendor := range r.vendors {
		result[i] = vendor
		i++
	}
	return result, nil
}
func (r *queryResolver) Concentrates(ctx context.Context) ([]model.Concentrate, error) {
	result := make([]model.Concentrate, len(r.vendors))
	i := 0
	for _, concentrate := range r.concentrates {
		result[i] = concentrate
		i++
	}
	return result, nil
}

type vendorResolver struct{ *Resolver }

func (r *vendorResolver) Urls(ctx context.Context, obj *model.Vendor) ([]*model.URL, error) {
	vendor := r.getVendor(obj.ID)
	if vendor == nil {
		return nil, fmt.Errorf("vendor not found")
	}
	result := make([]*model.URL, len(vendor.URLIDs))
	i := 0
	for _, urlID := range vendor.URLIDs {
		result[i] = r.getURL(*urlID)
		i++
	}
	return result, nil
}
