import {
  usePostUsers,
  usePostUsersLogin,
  usePostUsersLogout,
} from "@/api/endpoints/users/users";
import { useCallback } from "react";
import { useRouter } from "next/router";

interface LoginParams {
  email: string;
  password: string;
}

interface RegisterParams {
  name: string;
  email: string;
  password: string;
}

type UseAuthParams = {
  onLoginError?: () => void;
  onLogoutError?: () => void;
  onRegisterError?: () => void;
};

type UseAuth = (params?: UseAuthParams) => {
  login: (params: LoginParams) => void;
  logout: () => void;
  registerAndLogin: (params: RegisterParams) => void;
};

export const useAuth: UseAuth = (params) => {
  const router = useRouter();

  const loginMutation = usePostUsersLogin({
    mutation: {
      onSuccess: () => {
        router.push("/");
        console.log("onsuccess login");
      },
      onError: () => {
        params?.onLoginError?.();
      },
    },
    request: {
      withCredentials: true,
    },
  });

  const logoutMutation = usePostUsersLogout({
    mutation: {
      onSuccess: () => {
        router.push("/");
        console.log("onsuccess logout");
      },
      onError: () => {
        params?.onLogoutError?.();
      },
    },
    request: {
      withCredentials: true,
    },
  });

  const registerMutation = usePostUsers({
    mutation: {
      onError: () => {
        params?.onRegisterError?.();
      },
    },
    request: {
      withCredentials: true,
    },
  });

  const login = useCallback(
    ({ email, password }: LoginParams) => {
      loginMutation.mutate({ data: { email, password } });
    },
    [loginMutation]
  );

  const logout = useCallback(() => {
    logoutMutation.mutate();
  }, [logoutMutation]);

  const registerAndLogin = useCallback(
    async ({ email, password, name }: RegisterParams) => {
      await registerMutation
        .mutateAsync({
          data: {
            email,
            password,
            name,
          },
        })
        .then(() => login({ email, password }))
        .catch((error) => {
          console.error(error);
        });
    },
    [registerMutation, login]
  );

  return { login, logout, registerAndLogin };
};
