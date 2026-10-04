import type { CookieRef } from "nuxt/app";
import type { Ref } from "vue";
import type { PublicApi } from "~~/lib/api/public";
import type { UserOut } from "~~/lib/api/types/data-contracts";
import type { UserClient } from "~~/lib/api/user";

export interface IAuthContext {
  get token(): boolean | null;
  get attachmentToken(): string | null;
  get authToken(): string | null;

  /**
   * The current user object for the session. This is undefined if the session is not authorized.
   */
  user?: UserOut;

  /**
   * Returns true if the session is authorized.
   */
  isAuthorized(): boolean;

  /**
   * Invalidates the session by removing the token and the expiresAt.
   */
  invalidateSession(): void;

  /**
   * Logs out the user and calls the invalidateSession method.
   */
  logout(api: UserClient): ReturnType<UserClient["user"]["logout"]>;

  /**
   * Logs in the user and sets the authorization context via cookies
   */
  login(api: PublicApi, email: string, password: string, stayLoggedIn: boolean): ReturnType<PublicApi["login"]>;
}

class AuthContext implements IAuthContext {
  // eslint-disable-next-line no-use-before-define
  private static _instance?: AuthContext;

  // These names are part of the backend authentication protocol and must stay in sync with the API.
  private static readonly cookieTokenKey = "catalog.auth.session";
  private static readonly cookieAttachmentTokenKey = "catalog.auth.attachment_token";

  private _user: Ref<UserOut | undefined>;
  private _token: CookieRef<string | null>;
  private _attachmentToken: CookieRef<string | null>;
  private _authToken: Ref<string | null>;

  get user() {
    return this._user.value;
  }

  set user(value: UserOut | undefined) {
    this._user.value = value;
  }

  get token() {
    // @ts-ignore sometimes it's a boolean I guess?
    return this._token.value === "true" || this._token.value === true;
  }

  get attachmentToken() {
    return this._attachmentToken.value;
  }

  get authToken() {
    return this._authToken.value;
  }

  private constructor(token: string, attachmentToken: string) {
    this._user = useState<UserOut | undefined>("catalog.auth.user", () => undefined);
    this._authToken = useState<string | null>("catalog.auth.bearer", () => null);
    this._token = useCookie(token);
    this._attachmentToken = useCookie(attachmentToken);
  }

  static get instance() {
    if (!this._instance) {
      this._instance = new AuthContext(AuthContext.cookieTokenKey, AuthContext.cookieAttachmentTokenKey);
    }

    return this._instance;
  }

  isExpired() {
    return !this.token;
  }

  isAuthorized() {
    console.debug("isAuthorized", this.token);
    return this.token;
  }

  invalidateSession() {
    this.user = undefined;

    // Delete the cookies
    this._token.value = null;
    this._attachmentToken.value = null;
    this._authToken.value = null;
    console.log("Session invalidated");
  }

  async login(api: PublicApi, email: string, password: string, stayLoggedIn: boolean) {
    const r = await api.login(email, password, stayLoggedIn);

    if (!r.error) {
      const expiresAt = new Date(r.data.expiresAt);
      this._token = useCookie(AuthContext.cookieTokenKey);
      this._token.value = "true";
      this._authToken.value = r.data.token;
      this._attachmentToken = useCookie(AuthContext.cookieAttachmentTokenKey, {
        expires: expiresAt,
      });
      this._attachmentToken.value = r.data.attachmentToken;
    }

    return r;
  }

  async logout(api: UserClient) {
    const r = await api.user.logout();

    if (!r.error) {
      this.invalidateSession();
    }

    return r;
  }
}

export function useAuthContext(): IAuthContext {
  return AuthContext.instance;
}
