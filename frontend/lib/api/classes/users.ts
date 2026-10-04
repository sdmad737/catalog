import { BaseAPI, route } from "../base";
import type { ChangePassword, UserOut } from "../types/data-contracts";
import type { Result } from "../types/non-generated";

export class UserApi extends BaseAPI {
  public self() {
    return this.http.get<Result<UserOut>>({ url: route("/users/self") });
  }

  public logout() {
    return this.http.post<object, void>({ url: route("/users/logout") });
  }

  public delete() {
    return this.http.delete<void>({ url: route("/users/self") });
  }

  public changePassword(current: string, newPassword: string) {
    return this.http.put<ChangePassword, void>({
      url: route("/users/self/change-password"),
      body: {
        current,
        new: newPassword,
      },
    });
  }

  public list() {
    return this.http.get<{ items: UserOut[] }>({ url: route("/users") });
  }

  public setRole(id: string, role: "admin" | "staff" | "viewer") {
    return this.http.put<{ role: string }, Result<UserOut>>({
      url: route(`/users/${id}/role`),
      body: { role },
    });
  }

  public setDisabled(id: string, disabled: boolean) {
    return this.http.put<{ disabled: boolean }, Result<UserOut>>({
      url: route(`/users/${id}/status`),
      body: { disabled },
    });
  }

  public transferOwnership(userId: string) {
    return this.http.post<{ userId: string }, void>({
      url: route("/users/transfer-ownership"),
      body: { userId },
    });
  }
}
