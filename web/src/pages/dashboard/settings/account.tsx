import { AccountSettings } from '@/page-components/Dashboard/AccountSettings/AccountSettings'
import { Page } from '@/types/page'

const AccountPage: Page = () => {
  return <AccountSettings />
}

AccountPage.auth = true
export default AccountPage
