import { BaseAPI, route } from "../base";
import type { ItemOut, PersonCreate, PersonOut, PersonUpdate } from "../types/data-contracts";

export class PeopleApi extends BaseAPI {
  list(search = "") {
    return this.http.get<PersonOut[]>({
      url: route("/people", search ? { q: search } : {}),
    });
  }

  create(input: PersonCreate) {
    return this.http.post<PersonCreate, PersonOut>({
      url: route("/people"),
      body: input,
    });
  }

  update(id: string, input: PersonUpdate) {
    return this.http.put<PersonUpdate, PersonOut>({
      url: route(`/people/${id}`),
      body: input,
    });
  }

  checkout(
    itemId: string,
    input: {
      personId: string;
      dueAt?: string;
      condition: string;
      note?: string;
    }
  ) {
    return this.http.post<typeof input, ItemOut>({
      url: route(`/items/${itemId}/checkout`),
      body: input,
    });
  }

  returnItem(itemId: string, input: { locationId?: string; condition: string; note?: string }) {
    return this.http.post<typeof input, ItemOut>({
      url: route(`/items/${itemId}/return`),
      body: input,
    });
  }
}
