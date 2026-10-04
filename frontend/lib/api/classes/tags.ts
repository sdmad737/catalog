import { BaseAPI, route } from "../base";
import type { ItemOut } from "../types/data-contracts";

export type ItemTag = {
  id: string;
  token: string;
  createdAt: string;
  revokedAt?: string;
  verifiedAt?: string;
  lastScannedAt?: string;
};

export type PublicItemTag = {
  status: "active" | "revoked";
  name?: string;
  assetId?: string;
  organization?: string;
  itemStatus?: "available" | "missing";
};

export type ResolvedItemTag = { tag: ItemTag; item: ItemOut; canOperate: boolean };

export class TagsApi extends BaseAPI {
  list(itemId: string) {
    return this.http.get<{ items: ItemTag[] }>({
      url: route(`/items/${itemId}/tags`),
    });
  }

  assign(itemId: string) {
    return this.http.post<void, ItemTag>({
      url: route(`/items/${itemId}/tags`),
    });
  }

  replace(itemId: string) {
    return this.http.post<void, ItemTag>({
      url: route(`/items/${itemId}/tags/replace`),
    });
  }

  revoke(itemId: string, tagId: string) {
    return this.http.delete<void>({
      url: route(`/items/${itemId}/tags/${tagId}`),
    });
  }

  resolve(token: string) {
    return this.http.get<ResolvedItemTag>({
      url: route(`/tags/${encodeURIComponent(token)}/item`),
    });
  }
}
