package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	product_domain "github.com/ot07/next-bazaar/api/domain/product"
	db "github.com/ot07/next-bazaar/db/sqlc"
	"github.com/ot07/next-bazaar/util"
	"github.com/stretchr/testify/require"
)

func createSeed(t *testing.T, store db.Store, user db.User, category product_domain.Category, product product_domain.Product) (productId string) {
	ctx := context.Background()

	createdUser, err := store.CreateUser(ctx, db.CreateUserParams{
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	})
	require.NoError(t, err)

	createdCategory, err := store.CreateCategory(ctx, category.Name)
	require.NoError(t, err)

	createdProduct, err := store.CreateProduct(ctx, db.CreateProductParams{
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		CategoryID:    createdCategory.ID,
		SellerID:      createdUser.ID,
		ImageUrl:      product.ImageUrl,
	})
	require.NoError(t, err)

	return createdProduct.ID.String()
}

func createSeedDummy(t *testing.T, store db.Store, user db.User, category product_domain.Category, product product_domain.Product) (productId string) {
	return util.RandomUUID().String()
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	user, _ := randomUser(t)
	category := randomCategory(t)
	product := randomProduct(t, user, category)

	testCases := []struct {
		name          string
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store, user db.User, category product_domain.Category, product product_domain.Product) (productID string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeed,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				requireBodyMatchProduct(t, response.Body, product)
			},
		},
		{
			name: "NotFound",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: createSeedDummy,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name: "InvalidID",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: func(t *testing.T, store db.Store, user db.User, category product_domain.Category, product product_domain.Product) (productID string) {
				return "InvalidID"
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalError",
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Return(db.Product{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: createSeedDummy,
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store, cleanupStore := tc.buildStore(t)
			defer cleanupStore()

			productID := tc.createSeed(t, store, user, category, product)

			server := newTestServer(t, store)

			url := fmt.Sprintf("/api/v1/products/%s", productID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})
	}
}

func TestListProducts(t *testing.T) {
	t.Parallel()

	n := 5

	users := make([]db.User, n)
	categories := make([]product_domain.Category, n)
	products := make([]product_domain.Product, n)
	for i := 0; i < n; i++ {
		users[i], _ = randomUser(t)
		categories[i] = randomCategory(t)
		products[i] = randomProduct(t, users[i], categories[i])
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStore    func(t *testing.T) (store db.Store, cleanup func())
		createSeed    func(t *testing.T, store db.Store, users []db.User, categories []product_domain.Category, products []product_domain.Product) (productIDs []string)
		checkResponse func(t *testing.T, response *http.Response)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: func(t *testing.T, store db.Store, users []db.User, categories []product_domain.Category, products []product_domain.Product) (productIDs []string) {
				productIDs = make([]string, n)
				for i := 0; i < n; i++ {
					productIDs[i] = createSeed(t, store, users[i], categories[i], products[i])
				}
				return
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusOK, response.StatusCode)
				checkListProductsResponse(t, response.Body, products, 1, int32(n), 1, int64(n))
			},
		},
		{
			name: "PageIDNotFound",
			query: Query{
				pageSize: n,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: func(t *testing.T, store db.Store, users []db.User, categories []product_domain.Category, products []product_domain.Product) (productIDs []string) {
				productIDs = make([]string, n)
				for i := 0; i < n; i++ {
					productIDs[i] = createSeed(t, store, users[i], categories[i], products[i])
				}
				return
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageIDLessThanLowerLimit",
			query: Query{
				pageID:   0,
				pageSize: n,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: func(t *testing.T, store db.Store, users []db.User, categories []product_domain.Category, products []product_domain.Product) (productIDs []string) {
				productIDs = make([]string, n)
				for i := 0; i < n; i++ {
					productIDs[i] = createSeed(t, store, users[i], categories[i], products[i])
				}
				return
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeNotFound",
			query: Query{
				pageID: 1,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: func(t *testing.T, store db.Store, users []db.User, categories []product_domain.Category, products []product_domain.Product) (productIDs []string) {
				productIDs = make([]string, n)
				for i := 0; i < n; i++ {
					productIDs[i] = createSeed(t, store, users[i], categories[i], products[i])
				}
				return
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeLessThanLowerLimit",
			query: Query{
				pageID:   1,
				pageSize: 0,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: func(t *testing.T, store db.Store, users []db.User, categories []product_domain.Category, products []product_domain.Product) (productIDs []string) {
				productIDs = make([]string, n)
				for i := 0; i < n; i++ {
					productIDs[i] = createSeed(t, store, users[i], categories[i], products[i])
				}
				return
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PageSizeMoreThanUpperLimit",
			query: Query{
				pageID:   1,
				pageSize: 101,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				return newTestDBStore(t)
			},
			createSeed: func(t *testing.T, store db.Store, users []db.User, categories []product_domain.Category, products []product_domain.Product) (productIDs []string) {
				productIDs = make([]string, n)
				for i := 0; i < n; i++ {
					productIDs[i] = createSeed(t, store, users[i], categories[i], products[i])
				}
				return
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				pageID:   1,
				pageSize: 1,
			},
			buildStore: func(t *testing.T) (store db.Store, cleanup func()) {
				mockStore, cleanup := newMockStore(t)

				mockStore.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Return([]db.Product{}, sql.ErrConnDone)

				return mockStore, cleanup
			},
			createSeed: func(t *testing.T, store db.Store, users []db.User, categories []product_domain.Category, products []product_domain.Product) (productIDs []string) {
				return make([]string, n)
			},
			checkResponse: func(t *testing.T, response *http.Response) {
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store, cleanupStore := tc.buildStore(t)
			defer cleanupStore()

			tc.createSeed(t, store, users, categories, products)

			server := newTestServer(t, store)

			url := "/api/v1/products"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			response, err := server.app.Test(request, int(time.Second.Milliseconds()))
			require.NoError(t, err)

			tc.checkResponse(t, response)
		})

	}
}

func randomProduct(t *testing.T, user db.User, category product_domain.Category) product_domain.Product {
	price, err := util.RandomMoney()
	require.NoError(t, err)

	return product_domain.Product{
		Name:          util.RandomName(),
		Description:   sql.NullString{String: util.RandomString(30), Valid: true},
		Price:         price,
		StockQuantity: util.RandomInt32(10),
		CategoryID:    category.ID,
		Category:      category.Name,
		SellerID:      user.ID,
		Seller:        user.Name,
		ImageUrl:      sql.NullString{String: util.RandomImageUrl(), Valid: true},
	}
}

func randomCategory(t *testing.T) product_domain.Category {
	return product_domain.Category{
		Name: util.RandomName(),
	}
}

func requireBodyMatchProduct(t *testing.T, body io.ReadCloser, product product_domain.Product) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotProduct product_domain.ProductResponse
	err = json.Unmarshal(data, &gotProduct)
	require.NoError(t, err)

	require.NotEmpty(t, gotProduct.ID)
	requireProductResponseMatchProduct(t, gotProduct, product)

	err = body.Close()
	require.NoError(t, err)
}

func requireProductResponseMatchProduct(t *testing.T, gotProduct product_domain.ProductResponse, product product_domain.Product) {
	require.Equal(t, product.Name, gotProduct.Name)
	require.Equal(t, product.Description, gotProduct.Description.NullString)
	require.Equal(t, product.Price, gotProduct.Price)
	require.Equal(t, product.StockQuantity, gotProduct.StockQuantity)
	require.Equal(t, product.Category, gotProduct.Category)
	require.Equal(t, product.Seller, gotProduct.Seller)
	require.Equal(t, product.ImageUrl, gotProduct.ImageUrl.NullString)
}

func checkListProductsResponse(t *testing.T, body io.ReadCloser, products []product_domain.Product, pageID int32, pageSize int32, pageCount int64, totalCount int64) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse product_domain.ListProductsResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)

	require.Equal(t, pageID, gotResponse.Meta.PageID)
	require.Equal(t, pageSize, gotResponse.Meta.PageSize)
	require.Equal(t, pageCount, gotResponse.Meta.PageCount)
	require.Equal(t, totalCount, gotResponse.Meta.TotalCount)

	gotProducts := gotResponse.Data
	require.Len(t, gotProducts, len(products))

	sortedGotProducts := sortProductResponseByName(gotProducts)
	sortedProducts := sortProductByName(products)
	for i := range products {
		requireProductResponseMatchProduct(t, sortedGotProducts[i], sortedProducts[i])
	}

	err = body.Close()
	require.NoError(t, err)
}

func sortProductByName(products []product_domain.Product) []product_domain.Product {
	sortedProducts := make([]product_domain.Product, len(products))
	copy(sortedProducts, products)
	sort.Slice(sortedProducts, func(i, j int) bool {
		return sortedProducts[i].Name < sortedProducts[j].Name
	})
	return sortedProducts
}

func sortProductResponseByName(products product_domain.ProductsResponse) product_domain.ProductsResponse {
	sortedProducts := make([]product_domain.ProductResponse, len(products))
	copy(sortedProducts, products)
	sort.Slice(sortedProducts, func(i, j int) bool {
		return sortedProducts[i].Name < sortedProducts[j].Name
	})
	return sortedProducts
}
