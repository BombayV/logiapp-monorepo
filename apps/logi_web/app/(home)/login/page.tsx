"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    try {
      const response = await fetch(
        "https://poised-lexy-bombay-97849715.koyeb.app/auth/login",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ email, password }),
        }
      );

      if (response.ok) {
        const data = await response.json();
        document.cookie = `session=${data.token}; path=/;`;
        router.push("/dashboard");
      } else {
        const errorData = await response.json();
        setError(errorData.message || "Invalid credentials");
      }
    } catch (error) {
      setError("An unexpected error occurred.");
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <Card className="w-full max-w-sm">
        <form onSubmit={handleSubmit}>
          <CardHeader>
            <CardTitle className="text-2xl">Iniciar sesión</CardTitle>
            <CardDescription>
              Ingresa su email y contraseña para acceder a su cuenta.
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="me@gmail.com"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="password">Contraseña</Label>
              <Input
                id="password"
                type="password"
                required
                placeholder="********"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            {error && <p className="text-red-500 text-sm">{error}</p>}
          </CardContent>
          <CardFooter className="flex flex-col mt-4">
            <Button type="submit" className="w-full">
              Entrar
            </Button>
            <p className="mt-4 text-xs text-center text-muted-foreground">
              Si no tienes una cuenta, puedes registrarte mandando un correo a
              <Link href="mailto:me@bombayv.com" className="ml-1">
                <Button variant="link" className="px-0 text-xs">
                  nuestro email
                </Button>
              </Link>
              .
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
