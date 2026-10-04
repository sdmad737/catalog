<script setup lang="ts">
  import type { UserOut } from "~~/lib/api/types/data-contracts";
  import MdiAccountPlusOutline from "~icons/mdi/account-plus-outline";
  import MdiCrownOutline from "~icons/mdi/crown-outline";
  import MdiEmailOutline from "~icons/mdi/email-outline";

  definePageMeta({ middleware: ["auth"] });
  useHead({ title: "Staff — Catalog" });

  const api = useUserApi();
  const auth = useAuthContext();
  const toast = useNotifier();
  const inviteOpen = ref(false);
  const inviteLink = ref("");
  const invite = reactive({
    email: "",
    role: "staff" as "admin" | "staff" | "viewer",
  });

  const canManage = computed(() => auth.user?.role === "owner" || auth.user?.role === "admin");
  const isOwner = computed(() => auth.user?.role === "owner");

  const { data: staff, refresh } = useAsyncData("organization-staff", async () => {
    if (!canManage.value) {
      await navigateTo("/home");
      return [];
    }
    const response = await api.user.list();
    if (response.error) {
      toast.error("Could not load staff accounts.");
      return [];
    }
    return response.data.items;
  });

  async function createInvitation() {
    const expiresAt = new Date();
    expiresAt.setDate(expiresAt.getDate() + 7);
    const response = await api.group.createInvitation({
      email: invite.email || undefined,
      role: invite.role,
      uses: 1,
      expiresAt,
    });
    if (response.error) return toast.error("Could not create the invitation.");
    inviteLink.value = `${window.location.origin}?token=${response.data.token}`;
    toast.success("Invitation link created.");
  }

  async function copyInvite() {
    await navigator.clipboard.writeText(inviteLink.value);
    toast.success("Invitation link copied.");
  }

  async function changeRole(member: UserOut, event: Event) {
    const role = (event.target as HTMLSelectElement).value as "admin" | "staff" | "viewer";
    const response = await api.user.setRole(member.id, role);
    if (response.error) return toast.error("Could not change that role.");
    await refresh();
    toast.success(`${member.name}'s role is now ${role}.`);
  }

  async function toggleDisabled(member: UserOut) {
    const disabling = member.status !== "disabled";
    if (disabling && !window.confirm(`Disable ${member.name}? They will be signed out and unable to access Catalog.`))
      return;
    const response = await api.user.setDisabled(member.id, disabling);
    if (response.error) return toast.error("Could not update that account.");
    await refresh();
    toast.success(disabling ? "Account disabled." : "Account restored.");
  }

  async function transferOwnership(member: UserOut) {
    if (!window.confirm(`Transfer ownership to ${member.name}? You will become an administrator.`)) return;
    const response = await api.user.transferOwnership(member.id);
    if (response.error) return toast.error("Ownership could not be transferred.");
    toast.success(`Ownership transferred to ${member.name}.`);
    window.location.reload();
  }
</script>

<template>
  <BaseContainer class="catalog-page max-w-none">
    <AppPageHeader
      eyebrow="Organization"
      title="Staff & access"
      description="Manage the people who can sign in and what they are allowed to do."
    />

    <div
      class="mb-5 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-base-content/10 bg-base-100 p-4"
    >
      <div>
        <p class="font-extrabold">{{ auth.user?.groupName }}</p>
        <p class="mt-1 text-sm text-base-content/55">Owner, administrators, staff, and read-only viewers.</p>
      </div>
      <BaseButton class="btn-primary btn-sm" @click="inviteOpen = !inviteOpen"
        ><template #icon><MdiAccountPlusOutline /></template>Invite staff</BaseButton
      >
    </div>

    <BaseCard v-if="inviteOpen" class="mb-5">
      <template #title>Invite staff member</template>
      <div class="grid gap-4 border-t border-base-content/10 p-5 sm:grid-cols-[1fr_180px_auto] sm:items-end">
        <FormTextField v-model="invite.email" label="Email (optional restriction)" placeholder="person@example.com" />
        <label class="form-control"
          ><span class="label-text mb-2 text-sm font-bold">Role</span
          ><select v-model="invite.role" class="select select-bordered">
            <option value="admin">Administrator</option>
            <option value="staff">Staff</option>
            <option value="viewer">Viewer</option>
          </select></label
        >
        <BaseButton class="btn-primary" @click="createInvitation">Create invitation</BaseButton>
      </div>
      <div v-if="inviteLink" class="mx-5 mb-5 rounded-xl bg-base-200 p-4">
        <p class="text-xs font-bold uppercase tracking-wide text-base-content/45">Private invitation link</p>
        <p class="mt-2 break-all font-mono text-xs">{{ inviteLink }}</p>
        <button class="btn btn-secondary btn-sm mt-3" @click="copyInvite">
          <MdiEmailOutline />Copy invitation link
        </button>
        <p class="mt-2 text-xs text-base-content/50">
          Catalog does not send email automatically. Share this link privately; it expires in seven days and can be used
          once.
        </p>
      </div>
    </BaseCard>

    <BaseCard>
      <template #title>People with Catalog access</template>
      <div class="overflow-x-auto border-t border-base-content/10">
        <table class="table">
          <thead>
            <tr>
              <th>Person</th>
              <th>Role</th>
              <th>Status</th>
              <th>Joined</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="member in staff" :key="member.id">
              <td>
                <div class="font-extrabold">{{ member.name }}</div>
                <div class="text-xs text-base-content/50">
                  {{ member.email }}
                </div>
              </td>
              <td>
                <span v-if="member.role === 'owner'" class="badge badge-primary gap-1"><MdiCrownOutline />Owner</span
                ><select
                  v-else-if="isOwner"
                  class="select select-bordered select-sm"
                  :value="member.role"
                  @change="changeRole(member, $event)"
                >
                  <option value="admin">Administrator</option>
                  <option value="staff">Staff</option>
                  <option value="viewer">Viewer</option></select
                ><span v-else class="capitalize">{{ member.role }}</span>
              </td>
              <td>
                <span class="badge" :class="member.status === 'disabled' ? 'badge-neutral' : 'badge-success'">{{
                  member.status === "disabled" ? "Disabled" : "Active"
                }}</span>
              </td>
              <td><DateTime :date="member.joinedAt" /></td>
              <td class="text-right">
                <template v-if="member.id !== auth.user?.id && member.role !== 'owner'"
                  ><button class="btn btn-ghost btn-xs" @click="toggleDisabled(member)">
                    {{ member.status === "disabled" ? "Restore" : "Disable" }}</button
                  ><button
                    v-if="isOwner && member.status !== 'disabled'"
                    class="btn btn-ghost btn-xs"
                    @click="transferOwnership(member)"
                  >
                    Transfer ownership
                  </button></template
                ><span v-else class="text-xs text-base-content/35">Current account</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </BaseCard>
  </BaseContainer>
</template>
