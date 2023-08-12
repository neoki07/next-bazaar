import { NewProduct } from '@/page-components/Dashboard/Products'
import { Page } from '@/types/page'

const NewPage: Page = () => {
  return <NewProduct />
}

NewPage.auth = true
export default NewPage
