import { Products } from '@/page-components/Dashboard/Products'
import { Page } from '@/types/page'

const ProductsPage: Page = () => {
  return <Products />
}

ProductsPage.auth = true
export default ProductsPage
