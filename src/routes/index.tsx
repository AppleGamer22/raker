import { createFileRoute } from "@tanstack/react-router";

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
	return (
		<form
			onSubmit={(e) => {
				e.preventDefault();
			}}
		>
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Sign-in</FieldLegend>
					<FieldGroup>
						<Field>
							<FieldLabel>username</FieldLabel>
							<Input placeholder="username" />
						</Field>
						<Field>
							<FieldLabel>password</FieldLabel>
							<Input placeholder="password" type="password" />
						</Field>
						<Field orientation="horizontal">
							<Button type="submit">Sign-in</Button>
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
