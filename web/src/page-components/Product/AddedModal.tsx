import { CartProductList, useCart } from '@/features/cart'
import {
  Button,
  Center,
  Group,
  Modal,
  Stack,
  Text,
  UnstyledButton,
  createStyles,
} from '@mantine/core'
import { useRouter } from 'next/router'
import { useCallback, useEffect } from 'react'

const IMAGE_SIZE = 120

const useStyles = createStyles((theme) => ({
  title: {
    fontWeight: 700,
    fontSize: theme.fontSizes.lg,
  },
  continueButton: {
    color: theme.colors.gray[6],

    '&:hover': {
      color: theme.colors.gray[7],
    },
  },
}))

interface AddedModalProps {
  opened: boolean
  onClose: () => void
}

export function AddedModal({ opened, onClose }: AddedModalProps) {
  const { classes } = useStyles()
  const { data: cart, isFetching, refetch: refetchCart } = useCart()
  const router = useRouter()

  const viewCart = useCallback(() => {
    router.push('/cart')
  }, [router])

  useEffect(() => {
    if (opened) {
      refetchCart()
    }
  }, [opened, refetchCart])

  return (
    <Modal
      classNames={{ title: classes.title }}
      opened={opened}
      onClose={onClose}
      title="Added product to cart"
      centered
    >
      <Stack>
        <CartProductList
          cartProducts={cart?.products}
          isLoading={isFetching}
          imageSize={IMAGE_SIZE}
          editable={false}
        />
        <Group grow>
          <Button color="dark" onClick={viewCart}>
            View Cart
          </Button>

          <Center>
            <UnstyledButton onClick={onClose}>
              <Text className={classes.continueButton} size="sm" fw={500}>
                Continue shoppping
              </Text>
            </UnstyledButton>
          </Center>
        </Group>
      </Stack>
    </Modal>
  )
}
