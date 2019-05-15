package api

import (
	"context"
	"fmt"
	"log"

	"github.com/pscn/go4graphql/graph"
	"github.com/pscn/go4graphql/model"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	verbose      bool
	urls         map[string]*model.URL
	vendors      map[string]*model.Vendor
	concentrates map[string]*model.Concentrate
}

func NewResolver(verbose bool) *Resolver {
	r := &Resolver{
		verbose:      verbose,
		urls:         make(map[string]*model.URL, 10),
		vendors:      make(map[string]*model.Vendor, 10),
		concentrates: make(map[string]*model.Concentrate, 10),
	}
	// populate some data
	u := &model.URL{
		ID:          "http://www.capella.com",
		Description: "Homepage",
		URL:         "http://www.capella.com",
	}
	r.addURL(u)

	v := &model.Vendor{
		ID:   "CAP-Capella",
		Name: "Capella",
		Code: "CAP",
	}
	v.URLIDs = append(v.URLIDs, &u.ID)
	r.addVendor(v)

	c := &model.Concentrate{
		ID:   "CAP-Vanilla Custard",
		Name: "Vanilla Custard",
	}
	g := 1.012
	c.Gravity = &g
	c.URLIDs = append(v.URLIDs, &u.ID)
	r.addConcentrate(c)
	return r
}

// these would talk to the DB later

func (r *Resolver) addURL(url *model.URL) (*model.URL, error) {
	r.urls[url.ID] = url
	return url, nil
}

func (r *Resolver) getURL(id string) (*model.URL, error) {
	if result, ok := r.urls[id]; ok {
		return result, nil
	}
	return nil, nil
}

func (r *Resolver) getURLs(ids []*string) ([]*model.URL, error) {
	result := make([]*model.URL, len(ids))
	var err error
	for idx, id := range ids {
		result[idx], err = r.getURL(*id)
		if err != nil {
			return nil, nil
		}
	}
	return result, nil
}

func (r *Resolver) addVendor(vendor *model.Vendor) (*model.Vendor, error) {
	r.vendors[vendor.ID] = vendor
	return vendor, nil
}

func (r *Resolver) getVendor(id string) (*model.Vendor, error) {
	if result, ok := r.vendors[id]; ok {
		return result, nil
	}
	return nil, fmt.Errorf("vendor not found: %s", id)
}

func (r *Resolver) addConcentrate(concentrate *model.Concentrate) (*model.Concentrate, error) {
	r.concentrates[concentrate.ID] = concentrate
	return concentrate, nil
}

func (r *Resolver) getConcentrate(id string) (*model.Concentrate, error) {
	if result, ok := r.concentrates[id]; ok {
		return result, nil
	}
	return nil, fmt.Errorf("concentrate not found: %s", id)
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
	if r.verbose {
		log.Printf("Vendor: %+v", obj)
	}
	return r.getVendor(obj.VendorID)
}
func (r *concentrateResolver) Urls(ctx context.Context, obj *model.Concentrate) ([]*model.URL, error) {
	if r.verbose {
		log.Printf("Urls: %+v", obj)
	}
	return r.getURLs(obj.URLIDs)
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateVendor(ctx context.Context, input model.NewVendor) (*model.Vendor, error) {
	if r.verbose {
		log.Printf("CreateVendor: %+v", input)
	}
	id := fmt.Sprintf("%s-%s", input.Code, input.Name)
	return r.addVendor(&model.Vendor{
		ID:   id,
		Name: input.Name,
		Code: input.Code,
	})
}
func (r *mutationResolver) AddVendorURL(ctx context.Context, input model.NewVendorURL) (*model.URL, error) {
	if r.verbose {
		log.Printf("AddVendorURL: %+v", input)
	}
	vendor, err := r.getVendor(input.VendorID)
	if err != nil {
		return nil, err
	}
	url := &model.URL{
		ID:          input.URL,
		Description: input.Description,
		URL:         input.URL,
	}
	if _, err = r.addURL(url); err != nil {
		return nil, err
	}
	vendor.URLIDs = append(vendor.URLIDs, &input.URL)
	return url, nil
}
func (r *mutationResolver) CreateConcentrate(ctx context.Context, input model.NewConcentrate) (*model.Concentrate, error) {
	if r.verbose {
		log.Printf("CreateConcentrate: %+v", input)
	}
	vendor, err := r.getVendor(input.VendorID)
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("%s-%s", vendor.Code, input.Name)
	return r.addConcentrate(&model.Concentrate{
		ID:       id,
		Name:     input.Name,
		VendorID: vendor.ID,
	})
}
func (r *mutationResolver) AddConcentrateURL(ctx context.Context, input model.NewConcentrateURL) (*model.URL, error) {
	if r.verbose {
		log.Printf("AddConcentrateURL: %+v", input)
	}
	concentrate, err := r.getConcentrate(input.ConcentrateID)
	if err != nil {
		return nil, err
	}
	url := &model.URL{
		ID:          input.URL,
		Description: input.Description,
		URL:         input.URL,
	}
	if _, err = r.addURL(url); err != nil {
		return nil, err
	}
	concentrate.URLIDs = append(concentrate.URLIDs, &input.URL)
	return url, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Vendors(ctx context.Context) ([]model.Vendor, error) {
	if r.verbose {
		log.Printf("Vendors")
	}
	result := make([]model.Vendor, len(r.vendors))
	i := 0
	for _, vendor := range r.vendors {
		result[i] = *vendor
		i++
	}
	return result, nil
}
func (r *queryResolver) Concentrates(ctx context.Context) ([]model.Concentrate, error) {
	result := make([]model.Concentrate, len(r.vendors))
	i := 0
	for _, concentrate := range r.concentrates {
		result[i] = *concentrate
		i++
	}
	return result, nil
}

type vendorResolver struct{ *Resolver }

func (r *vendorResolver) Urls(ctx context.Context, obj *model.Vendor) ([]*model.URL, error) {
	if r.verbose {
		log.Printf("Urls: %+v", obj)
	}
	vendor, err := r.getVendor(obj.ID)
	if err != nil {
		return nil, err
	}
	result := make([]*model.URL, len(vendor.URLIDs))
	i := 0
	for _, urlID := range vendor.URLIDs {
		result[i], err = r.getURL(*urlID)
		if err != nil {
			return nil, err
		}
		i++
	}
	return result, nil
}
