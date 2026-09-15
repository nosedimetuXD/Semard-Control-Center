"use client";

import { useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { setAuthToken } from "@/lib/api";

function CallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get("token");
    const status = searchParams.get("status");

    if (token) {
      setAuthToken(token);
      router.replace("/?auth_success=true");
      return;
    }

    if (status === "NEEDS_REGISTRATION") {
      const email = searchParams.get("email") || "";
      const name = searchParams.get("name") || "";
      const googleId = searchParams.get("google_id") || "";
      router.replace(`/?status=NEEDS_REGISTRATION&email=${encodeURIComponent(email)}&name=${encodeURIComponent(name)}&google_id=${encodeURIComponent(googleId)}`);
      return;
    }

    if (status === "PENDING_APPROVAL") {
      const email = searchParams.get("email") || "";
      router.replace(`/?status=PENDING_APPROVAL&email=${encodeURIComponent(email)}`);
      return;
    }

    if (status === "REJECTED") {
      const msg = searchParams.get("msg") || "";
      router.replace(`/?status=REJECTED&msg=${encodeURIComponent(msg)}`);
      return;
    }

    router.replace("/");
  }, [router, searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-950 text-slate-200">
      <div className="bg-slate-900 border border-slate-800 p-8 rounded-xl max-w-sm text-center space-y-4 shadow-xl">
        <div className="w-10 h-10 border-4 border-teal-500/20 border-t-teal-400 rounded-full animate-spin mx-auto" />
        <h2 className="text-base font-semibold text-white">Procesando autenticación</h2>
        <p className="text-xs text-slate-400">Verificando credenciales con la Universidad de Cartagena...</p>
      </div>
    </div>
  );
}

export default function AuthCallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-slate-950 text-slate-200">
        <div className="w-8 h-8 border-4 border-teal-500/20 border-t-teal-400 rounded-full animate-spin" />
      </div>
    }>
      <CallbackContent />
    </Suspense>
  );
}
