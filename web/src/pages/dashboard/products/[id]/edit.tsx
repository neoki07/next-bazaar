import { EditProduct } from '@/page-components/Dashboard/Products/EditProduct'
import { Page } from '@/types/page'
import { useRouter } from 'next/router'

const EditPage: Page = () => {
  const router = useRouter()
  const { id } = router.query

  if (Array.isArray(id)) {
    throw new Error('id is array:' + JSON.stringify(id))
  }

  if (id === undefined) {
    return null
  }

  return <EditProduct productId={id} />
}

EditPage.auth = true
export default EditPage
