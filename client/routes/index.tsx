import { useMutation } from "@connectrpc/connect-query";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";

import { signInInstagram } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { Button } from "@/components/ui/button";
import { CardContent } from "@/components/ui/card";
import { Field, FieldGroup, FieldLabel, FieldLegend, FieldSet } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
// import { createRootRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({ component: AuthPage, ssr: false });

function SignUpForm() {
	return (
		<form
			onSubmit={(e) => {
				e.preventDefault();
			}}
		>
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Sign-up</FieldLegend>
					<FieldGroup>
						<Field>
							<FieldLabel>username</FieldLabel>
							<Input placeholder="username" />
						</Field>
						<Field>
							<FieldLabel>password</FieldLabel>
							<Input placeholder="password" type="password" />
						</Field>
						<Field>
							<FieldLabel>session ID</FieldLabel>
							<Input placeholder="session ID cookie value" />
						</Field>
						<Field>
							<FieldLabel>user ID</FieldLabel>
							<Input placeholder="user ID cookie value" />
						</Field>
						<Field orientation="horizontal">
							<Button type="submit">Sign-up</Button>
						</Field>
					</FieldGroup>
				</FieldSet>
			</FieldGroup>
		</form>
	);
}

function SignInForm() {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const signInMutation = useMutation(signInInstagram);

	return (
		<form
			onSubmit={(e) => {
				e.preventDefault();
				signInMutation.mutate({ username, password });
			}}
		>
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Sign-in</FieldLegend>
					<FieldGroup>
						<Field>
							<FieldLabel>username</FieldLabel>
							<Input
								autoComplete="username"
								placeholder="username"
								value={username}
								onChange={(e) => {
									setUsername(e.target.value);
								}}
							/>
						</Field>
						<Field>
							<FieldLabel>password</FieldLabel>
							<Input
								autoComplete="current-password"
								placeholder="password"
								type="password"
								value={password}
								onChange={(e) => {
									setPassword(e.target.value);
								}}
							/>
						</Field>
						{signInMutation.isError ? (
							<p className="text-sm text-destructive">{signInMutation.error.message}</p>
						) : null}
						{signInMutation.isSuccess ? <p className="text-sm text-green-600">Signed in.</p> : null}
						<Field orientation="horizontal">
							<Button disabled={signInMutation.isPending} type="submit">
								{signInMutation.isPending ? "Signing in..." : "Sign-in"}
							</Button>
						</Field>
					</FieldGroup>
				</FieldSet>
			</FieldGroup>
		</form>
	);
}

function SignedOut() {
	return (
		<>
			<CardContent>
				<SignUpForm />
				<Separator className="my-3" />
				<SignInForm />
			</CardContent>
		</>
	);
}

function UpdateForm() {
	return (
		<form
			onSubmit={(e) => {
				e.preventDefault();
			}}
		>
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Update</FieldLegend>
					<FieldGroup>
						<Field>
							<FieldLabel>password</FieldLabel>
							<Input placeholder="password" type="password" />
						</Field>
						<Field>
							<FieldLabel>session ID</FieldLabel>
							<Input placeholder="session ID cookie value" />
						</Field>
						<Field>
							<FieldLabel>user ID</FieldLabel>
							<Input placeholder="user ID cookie value" />
						</Field>
						<Field orientation="horizontal">
							<Button type="submit">Update</Button>
						</Field>
					</FieldGroup>
				</FieldSet>
			</FieldGroup>
		</form>
	);
}

// oxlint-disable-next-line no-unused-vars
function SignedIn() {
	return (
		<>
			<CardContent>
				<Separator className="my-3" />
				<UpdateForm />
			</CardContent>
		</>
	);
}

function AuthPage() {
	return <SignedOut />;
}
