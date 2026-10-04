<script setup lang="ts">
  import MdiAccountPlusOutline from "~icons/mdi/account-plus-outline";
  import MdiArrowRight from "~icons/mdi/arrow-right";
  import MdiCloudAlertOutline from "~icons/mdi/cloud-alert-outline";
  import MdiLockOutline from "~icons/mdi/lock-outline";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";
  import MdiPackageVariantClosed from "~icons/mdi/package-variant-closed";
  import MdiQrcodeScan from "~icons/mdi/qrcode-scan";
  import MdiWrenchClockOutline from "~icons/mdi/wrench-clock-outline";

  useHead({ title: "Catalog — Know what you have and where it is" });

  definePageMeta({
    layout: "empty",
    middleware: [() => (useAuthContext().isAuthorized() ? "/home" : undefined)],
  });

  const ctx = useAuthContext();
  const api = usePublicApi();
  const toast = useNotifier();
  const route = useRoute();
  const router = useRouter();

  const username = ref("");
  const email = ref("");
  const password = ref("");
  const loginPassword = ref("");
  const canRegister = ref(false);
  const remember = ref(false);
  const loading = ref(false);
  const registerForm = ref(false);

  const { data: status, pending: statusPending } = useAsyncData(async () => {
    const response = await api.status();
    return response.error ? null : response.data;
  });

  whenever(status, current => {
    if (current?.demo) {
      email.value = "demo@example.com";
      loginPassword.value = "demo";
    }
  });

  const groupToken = computed<string>({
    get: () => (typeof route.query.token === "string" ? route.query.token : ""),
    set: token => router.push({ query: token ? { token } : {} }),
  });

  onMounted(() => {
    if (groupToken.value) registerForm.value = true;
  });

  async function registerUser() {
    loading.value = true;
    const { error } = await api.register({
      name: username.value,
      email: email.value,
      password: password.value,
      token: groupToken.value,
    });
    loading.value = false;
    if (error) {
      toast.error("We couldn't create your account. Check the details and try again.");
      return;
    }
    toast.success("Account created. You can sign in now.");
    registerForm.value = false;
  }

  async function login() {
    loading.value = true;
    const { error } = await ctx.login(api, email.value, loginPassword.value, remember.value);
    loading.value = false;
    if (error) {
      toast.error("That email and password combination didn't work.");
      return;
    }
    toast.success("Welcome back.");
    const redirect = typeof route.query.redirect === "string" ? route.query.redirect : "";
    navigateTo(redirect.startsWith("/") && !redirect.startsWith("//") ? redirect : "/home");
  }

  const benefits = [
    {
      icon: MdiPackageVariantClosed,
      title: "One source of truth",
      text: "Keep instruments, accessories, and stage equipment together.",
    },
    {
      icon: MdiMapMarkerOutline,
      title: "Always know where it is",
      text: "Organize gear by room, case, shelf, or nested location.",
    },
    {
      icon: MdiQrcodeScan,
      title: "Find gear in seconds",
      text: "Use asset IDs and QR labels to jump straight to an item.",
    },
    {
      icon: MdiWrenchClockOutline,
      title: "Stay ahead of maintenance",
      text: "Keep service notes, documents, and reminders attached to each item.",
    },
  ];
</script>

<template>
  <div class="min-h-screen bg-base-200 lg:grid lg:grid-cols-[1.08fr_0.92fr]">
    <AppToast />

    <section
      class="catalog-grid-pattern relative hidden min-h-screen overflow-hidden bg-neutral px-12 py-10 text-neutral-content lg:flex lg:flex-col xl:px-20"
    >
      <div class="absolute -right-32 -top-32 h-96 w-96 rounded-full bg-primary/25 blur-3xl" />
      <div class="relative flex items-center gap-3">
        <div
          class="grid h-12 w-12 place-items-center rounded-2xl bg-primary text-primary-content shadow-xl shadow-primary/20"
        >
          <AppLogo class="h-8 w-8" />
        </div>
        <div>
          <p class="font-display text-2xl font-extrabold">Catalog</p>
          <p class="text-xs text-neutral-content/50">Inventory, without the clutter</p>
        </div>
      </div>

      <div class="relative my-auto max-w-2xl py-16">
        <p class="mb-5 text-xs font-bold uppercase tracking-[0.2em] text-accent">Built for equipment that matters</p>
        <h1 class="text-5xl font-extrabold leading-[1.08] xl:text-6xl">
          Know what you have.<br /><span class="text-primary">Know where it is.</span>
        </h1>
        <p class="mt-6 max-w-xl text-lg leading-8 text-neutral-content/65">
          A calmer way to organize instruments, production gear, documents, and maintenance history.
        </p>

        <div class="mt-12 grid gap-4 sm:grid-cols-2">
          <article
            v-for="benefit in benefits"
            :key="benefit.title"
            class="rounded-2xl border border-white/10 bg-white/5 p-5 backdrop-blur-sm"
          >
            <component :is="benefit.icon" class="h-6 w-6 text-accent" />
            <h2 class="mt-4 text-base font-bold">{{ benefit.title }}</h2>
            <p class="mt-1 text-sm leading-6 text-neutral-content/55">{{ benefit.text }}</p>
          </article>
        </div>
      </div>

      <div class="relative flex items-center justify-between text-xs text-neutral-content/40">
        <span>Designed by Shahin Madantschi</span><span v-if="status?.build">v{{ status.build.version }}</span>
      </div>
    </section>

    <main class="flex min-h-screen items-center justify-center px-4 py-10 sm:px-8">
      <div class="w-full max-w-md">
        <div class="mb-10 flex items-center justify-center gap-3 lg:hidden">
          <div class="grid h-11 w-11 place-items-center rounded-xl bg-primary text-primary-content">
            <AppLogo class="h-7 w-7" />
          </div>
          <span class="font-display text-2xl font-extrabold">Catalog</span>
        </div>

        <div class="mb-8">
          <p class="catalog-eyebrow">{{ registerForm ? "Create your workspace" : "Welcome back" }}</p>
          <h1 class="mt-2 text-3xl font-extrabold">{{ registerForm ? "Start organizing" : "Sign in to Catalog" }}</h1>
          <p class="catalog-muted mt-2 text-sm leading-6">
            {{
              registerForm
                ? "Create an account to begin building your inventory."
                : "Everything you track is right where you left it."
            }}
          </p>
        </div>

        <form v-if="registerForm" class="space-y-4" @submit.prevent="registerUser">
          <div v-if="groupToken" class="rounded-xl border border-primary/20 bg-primary/5 p-4 text-sm">
            <p class="font-semibold">You were invited to an existing workspace.</p>
            <button
              type="button"
              class="mt-1 text-xs font-semibold text-primary hover:underline"
              @click="groupToken = ''"
            >
              Create a separate workspace instead
            </button>
          </div>
          <FormTextField v-model="username" label="Your name" placeholder="Shahin Madantschi" />
          <FormTextField v-model="email" type="email" label="Email address" placeholder="you@example.com" />
          <FormPassword v-model="password" label="Create a password" />
          <PasswordScore v-model:valid="canRegister" :password="password" />
          <button
            type="submit"
            class="btn btn-primary mt-2 w-full"
            :class="{ loading }"
            :disabled="loading || !canRegister"
          >
            Create account <MdiArrowRight class="h-5 w-5" />
          </button>
        </form>

        <form v-else class="space-y-4" @submit.prevent="login">
          <div v-if="status?.demo" class="rounded-xl border border-info/20 bg-info/5 p-4 text-sm">
            <strong>Demo workspace:</strong> credentials are already filled in.
          </div>
          <FormTextField v-model="email" type="email" label="Email address" placeholder="you@example.com" />
          <FormPassword v-model="loginPassword" label="Password" />
          <FormCheckbox v-model="remember" label="Keep me signed in" />
          <button type="submit" class="btn btn-primary mt-2 w-full" :class="{ loading }" :disabled="loading">
            Sign in <MdiArrowRight class="h-5 w-5" />
          </button>
        </form>

        <div class="mt-6 text-center">
          <button
            v-if="status?.allowRegistration"
            class="inline-flex min-h-[44px] items-center gap-2 text-sm font-semibold text-primary hover:underline"
            @click="registerForm = !registerForm"
          >
            <MdiAccountPlusOutline v-if="!registerForm" class="h-5 w-5" />{{
              registerForm ? "Already have an account? Sign in" : "New to Catalog? Create an account"
            }}
          </button>
          <p v-else-if="statusPending" class="catalog-muted text-sm">Connecting to Catalog…</p>
          <p v-else-if="status" class="catalog-muted inline-flex items-center gap-2 text-sm">
            <MdiLockOutline class="h-4 w-4" />Registration is managed by your administrator
          </p>
          <p
            v-else
            class="inline-flex items-center gap-2 rounded-lg bg-error/10 px-3 py-2 text-sm font-semibold text-error"
          >
            <MdiCloudAlertOutline class="h-4 w-4" />Catalog service is unavailable
          </p>
        </div>

        <p class="catalog-muted mt-10 text-center text-xs lg:hidden">A student project by Shahin Madantschi</p>
      </div>
    </main>
  </div>
</template>
