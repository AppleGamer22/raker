import { useMutation, useQuery } from "@connectrpc/connect-query";
import {
	startAuthentication,
	startRegistration,
	type PublicKeyCredentialCreationOptionsJSON,
	type PublicKeyCredentialRequestOptionsJSON,
} from "@simplewebauthn/browser";
import { useForm } from "@tanstack/react-form";
import { createFileRoute } from "@tanstack/react-router";
import { PlusIcon, UserKeyIcon, XIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import z from "zod";

import type { Passkey } from "@/buf/raker/v1/passkey/passkey_pb";
import {
	signInInstagram,
	editCategory,
	editUserCredentials,
	beginSignUp,
	finishSignUp,
	beginSignIn,
	finishSignIn,
	renamePasskey,
	deletePasskey,
	getPasskeysList,
} from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { Button } from "@/components/ui/button";
import { CardContent } from "@/components/ui/card";
import { Field, FieldContent, FieldError, FieldGroup, FieldLabel, FieldLegend, FieldSet } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { InputGroup, InputGroupAddon, InputGroupButton, InputGroupInput } from "@/components/ui/input-group";
import { Progress } from "@/components/ui/progress";
import { Separator } from "@/components/ui/separator";
import { useConfirmationDialog } from "@/hooks/use-confirmation-dialog";
import { useUser } from "@/hooks/user-provider";
// import { createRootRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({ component: AuthPage, ssr: false });

function SignUpForm() {
	const beginSignUpMutation = useMutation(beginSignUp);
	const finishSignUpMutation = useMutation(finishSignUp);
	const [username, setUsername] = useState("");
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
							<Input
								placeholder="username"
								value={username}
								onChange={(e) => setUsername(e.target.value)}
							/>
						</Field>
						{/* <Field>
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
						</Field> */}
						<Field orientation="horizontal">
							<Button type="submit">Sign-up</Button>
							<Button
								type="submit"
								disabled={beginSignUpMutation.isPending || finishSignUpMutation.isPending}
								onClick={async () => {
									try {
										const beginRes = await beginSignUpMutation.mutateAsync({ username });
										const { publicKey } = JSON.parse(beginRes.optionsJson) as {
											publicKey: PublicKeyCredentialCreationOptionsJSON;
										};

										const attResp = await startRegistration({ optionsJSON: publicKey });
										await finishSignUpMutation.mutateAsync({
											sessionId: beginRes.sessionId,
											responseJson: JSON.stringify(attResp),
										});

										toast.success("Registered Passkey", {
											position: "top-center",
										});
										setUsername("");
									} catch (err) {
										console.error(err);
										toast.error((err as Error).message, {
											position: "top-center",
										});
									}
								}}
							>
								<UserKeyIcon className="h-4 w-4" /> Passkey Sign-up
							</Button>
						</Field>
					</FieldGroup>
				</FieldSet>
			</FieldGroup>
		</form>
	);
}

function PasskeysForm() {
	const { data: passkeys, error: passkeysListError, refetch } = useQuery(getPasskeysList, {});
	const { confirm, DialogComponent } = useConfirmationDialog();

	function PasskeyRow({ passkey }: { passkey: Passkey }) {
		const renamePasskeyMutation = useMutation(renamePasskey);
		const deletePasskeyMutation = useMutation(deletePasskey);
		const [name, setName] = useState(passkey.name);

		useEffect(() => {
			setName(passkey.name);
		}, [passkey.name]);

		const trimmedName = name.trim();
		const isRenameDisabled = !trimmedName || trimmedName === passkey.name || renamePasskeyMutation.isPending;

		return (
			<div className="flex flex-col gap-3 rounded-lg border bg-card p-3 md:flex-row md:items-center">
				<div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-md border bg-muted text-[10px] font-semibold tracking-wide text-muted-foreground uppercase">
					Logo
				</div>
				<div className="min-w-0 flex-1 space-y-2">
					<div className="flex items-center gap-2 text-sm">
						<span className="font-medium">Passkey name</span>
						<span className="truncate text-muted-foreground">{passkey.name}</span>
					</div>
					<Input
						value={name}
						onChange={(e) => setName(e.target.value)}
						placeholder="Rename this passkey"
						aria-label={`Rename passkey ${passkey.name}`}
					/>
				</div>
				<div className="flex flex-wrap items-center gap-2 md:justify-end">
					<Button
						type="button"
						variant="outline"
						disabled={isRenameDisabled}
						onClick={async () => {
							try {
								await renamePasskeyMutation.mutateAsync({
									id: passkey.id,
									name: trimmedName,
									aaguid: passkey.aaguid,
								});
								await refetch();
								toast.success("Renamed passkey", {
									position: "top-center",
								});
							} catch (err) {
								toast.error((err as Error).message, {
									position: "top-center",
								});
							}
						}}
					>
						Save name
					</Button>
					<Button
						type="button"
						variant="destructive"
						disabled={(passkeys?.passkeys ?? []).length < 2 || deletePasskeyMutation.isPending}
						onClick={async () => {
							const confirmed = await confirm({
								title: "Delete Passkey",
								description: `Are you sure you want to delete the passkey "${passkey.name}"? This action cannot be undone.`,
								confirmText: "Delete",
								cancelText: "Cancel",
								isDestructive: true,
							});

							if (!confirmed) {
								return;
							}

							try {
								await deletePasskeyMutation.mutateAsync({
									id: passkey.id,
									name: passkey.name,
									aaguid: passkey.aaguid,
								});
								await refetch();
								toast.success("Deleted passkey", {
									position: "top-center",
								});
							} catch (err) {
								toast.error((err as Error).message, {
									position: "top-center",
								});
							}
						}}
					>
						Delete
					</Button>
				</div>
			</div>
		);
	}

	const passkeyRows = passkeys?.passkeys ?? [];
	return (
		<>
			<form
				onSubmit={(e) => {
					e.preventDefault();
				}}
			>
				<FieldGroup>
					<FieldSet>
						<FieldLegend>Passkeys</FieldLegend>
						<FieldGroup>
							{passkeysListError ? <FieldError>{passkeysListError.message}</FieldError> : null}
							{passkeyRows.length ? (
								passkeyRows.map((passkey) => (
									<PasskeyRow key={passkey.id.toString()} passkey={passkey} />
								))
							) : (
								<p className="text-sm text-muted-foreground">No passkeys have been registered yet.</p>
							)}
						</FieldGroup>
					</FieldSet>
				</FieldGroup>
				{passkeys === undefined && <Progress value={null} className="pt-2" />}
			</form>
			<DialogComponent />
		</>
	);
}

function SignInForm() {
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const signInMutation = useMutation(signInInstagram);
	const beginSignInMutation = useMutation(beginSignIn);
	const finishSignInMutation = useMutation(finishSignIn);

	return (
		<form
			onSubmit={async (e) => {
				e.preventDefault();
				try {
					await signInMutation.mutateAsync({ username, password });
					location.reload();
				} catch (err) {
					toast.error((err as Error).message, {
						position: "top-center",
					});
				}
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
								onChange={(e) => setUsername(e.target.value)}
							/>
						</Field>
						<Field>
							<FieldLabel>password</FieldLabel>
							<Input
								autoComplete="current-password"
								placeholder="password"
								type="password"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
							/>
						</Field>
						<Field orientation="horizontal">
							<Button disabled={signInMutation.isPending} type="submit">
								{signInMutation.isPending ? "Signing in..." : "Sign-in"}
							</Button>
							<Button
								// type="submit"
								disabled={beginSignInMutation.isPending || finishSignInMutation.isPending}
								onClick={async () => {
									try {
										const beginRes = await beginSignInMutation.mutateAsync({ username: "" });
										const optionsJSON = JSON.parse(
											beginRes.optionsJson,
										) as PublicKeyCredentialRequestOptionsJSON;

										const attResp = await startAuthentication({ optionsJSON });
										await finishSignInMutation.mutateAsync({
											sessionId: beginRes.sessionId,
											responseJson: JSON.stringify(attResp),
										});
										location.reload();
									} catch (err) {
										console.error(err);
										toast.error((err as Error).message, {
											position: "top-center",
										});
									}
								}}
							>
								<UserKeyIcon className="h-4 w-4" /> Passkey Sign-in
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
	const editUserCredentialsMutation = useMutation(editUserCredentials);
	const [password, setPassword] = useState("");
	const [sessionID, setSessionID] = useState("");
	const [userID, setUserID] = useState("");

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
							<Input
								placeholder="new password"
								type="password"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
							/>
						</Field>
						<Field>
							<FieldLabel>session ID</FieldLabel>
							<Input
								placeholder="session ID cookie value"
								value={sessionID}
								onChange={(e) => setSessionID(e.target.value)}
							/>
						</Field>
						<Field>
							<FieldLabel>user ID</FieldLabel>
							<Input
								placeholder="user ID cookie value"
								value={userID}
								onChange={(e) => setUserID(e.target.value)}
							/>
						</Field>
						<Field orientation="horizontal">
							<Button
								type="submit"
								onClick={async () => {
									try {
										await editUserCredentialsMutation.mutateAsync({
											password,
											sessionId: sessionID,
											userId: userID,
										});
										setPassword("");
										setSessionID("");
										setUserID("");
										toast.success("Updated credentials", {
											position: "top-center",
										});
									} catch (err) {
										toast.error((err as Error).message, {
											position: "top-center",
										});
									}
								}}
							>
								Update
							</Button>
							<Button
								variant="destructive"
								onClick={async () => {
									try {
										await cookieStore.delete("jwt");
										location.reload();
									} catch (err) {
										toast.error((err as Error).message, {
											position: "top-center",
										});
									}
								}}
							>
								Sign-out
							</Button>
						</Field>
					</FieldGroup>
				</FieldSet>
			</FieldGroup>
			{editUserCredentialsMutation.isPending && <Progress value={null} className="pt-2" />}
		</form>
	);
}

const updateCategoriesSchema = z.object({
	categories: z.array(z.string().catch("")),
});

function Categories() {
	const { categories, categoriesError, isCategoriesPending, setShouldRefetchCategories } = useUser();
	const [newCategory, setNewCategory] = useState("");
	const form = useForm({
		defaultValues: {
			categories,
		},
		validators: {
			onBlur: updateCategoriesSchema,
		},
	});
	const editCategoryMutation = useMutation(editCategory);
	const { confirm, DialogComponent } = useConfirmationDialog();

	// Update form when categories are fetched
	useEffect(() => {
		form.setFieldValue("categories", categories);
	}, [categories, form]);

	return (
		<form
			onSubmit={(e) => {
				e.preventDefault();
			}}
		>
			<FieldLegend>Categories</FieldLegend>
			<FieldGroup>
				{categoriesError ? <FieldError>{categoriesError}</FieldError> : null}
				<form.Field name="categories" mode="array">
					{(field) => {
						const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
						return (
							<FieldSet>
								<FieldGroup>
									{field.state.value.map((_, i) => (
										<form.Field key={i} name={`categories[${i}]`}>
											{(subField) => {
												const isSubFieldInvalid =
													subField.state.meta.isTouched && !subField.state.meta.isValid;
												return (
													<Field orientation="horizontal" data-invalid={isSubFieldInvalid}>
														<FieldContent>
															<InputGroup>
																<InputGroupInput
																	name={subField.name}
																	value={subField.state.value}
																	onBlur={subField.handleBlur}
																	onChange={(e) =>
																		subField.handleChange(e.target.value)
																	}
																	aria-invalid={isSubFieldInvalid}
																	placeholder={
																		i < categories.length
																			? categories[i]
																			: "New Category Name"
																	}
																></InputGroupInput>

																<InputGroupAddon align="inline-end">
																	<InputGroupButton
																		type="button"
																		variant="ghost"
																		size="icon-xs"
																		onClick={async () => {
																			const confirmed = await confirm({
																				title: "Delete Category",
																				description: `Are you sure you want to delete the category "${subField.state.value}"? This action cannot be undone.`,
																				confirmText: "Delete",
																				cancelText: "Cancel",
																				isDestructive: true,
																			});

																			if (!confirmed) return;

																			try {
																				await editCategoryMutation.mutateAsync({
																					oldCategory: subField.state.value,
																					newCategory: "",
																				});
																				field.removeValue(i);
																				setShouldRefetchCategories(true);
																			} catch (err) {
																				toast.error((err as Error).message, {
																					position: "top-center",
																				});
																			}
																		}}
																		aria-label={`Remove Category ${i + 1}`}
																	>
																		<XIcon />
																	</InputGroupButton>
																</InputGroupAddon>
															</InputGroup>
															{isSubFieldInvalid && (
																<FieldError errors={subField.state.meta.errors} />
															)}
														</FieldContent>
													</Field>
												);
											}}
										</form.Field>
									))}
								</FieldGroup>
								<Field orientation="horizontal">
									<FieldContent>
										<InputGroup>
											<InputGroupInput
												value={newCategory}
												onChange={(e) => {
													setNewCategory(e.target.value);
												}}
												placeholder="New Category Name"
											/>
											<InputGroupAddon align="inline-end">
												<InputGroupButton
													type="button"
													variant="ghost"
													size="icon-xs"
													onClick={() => setNewCategory("")}
													aria-label={`Reset New Category Name`}
													disabled={!newCategory.trim()}
												>
													<XIcon />
													<span className="sr-only">Remove category</span>
												</InputGroupButton>
												<InputGroupButton
													type="button"
													onClick={async () => {
														const trimmedNewCategoryName = newCategory.trim();
														if (!newCategory) {
															return;
														} else if (field.state.value.includes(trimmedNewCategoryName)) {
															toast.error(
																<label>
																	Category name <b>{trimmedNewCategoryName}</b> is
																	already part of the categories list
																</label>,
																{
																	position: "top-center",
																},
															);
															return;
														}

														try {
															await editCategoryMutation.mutateAsync({
																oldCategory: trimmedNewCategoryName,
																newCategory: trimmedNewCategoryName,
															});

															setShouldRefetchCategories(true);
															field.pushValue(trimmedNewCategoryName);
															setNewCategory("");
														} catch (err) {
															toast.error((err as Error).message, {
																position: "top-center",
															});
														}
													}}
													disabled={!newCategory.trim()}
												>
													<PlusIcon />
													<span className="sr-only">Add category</span>
												</InputGroupButton>
											</InputGroupAddon>
										</InputGroup>
									</FieldContent>
								</Field>
								{isInvalid && <FieldError errors={field.state.meta.errors} />}
							</FieldSet>
						);
					}}
				</form.Field>
			</FieldGroup>
			{(isCategoriesPending || editCategoryMutation.isPending) && <Progress value={null} className="pt-2" />}
			<DialogComponent />
		</form>
	);
}

// oxlint-disable-next-line no-unused-vars
function SignedIn() {
	return (
		<>
			<CardContent>
				<Categories />
				<Separator className="my-3" />
				<PasskeysForm />
				<Separator className="my-3" />
				<UpdateForm />
			</CardContent>
		</>
	);
}

function AuthPage() {
	const { username, refetchCategoriesIfRequested } = useUser();

	useEffect(() => {
		return () => {
			refetchCategoriesIfRequested();
		};
	}, [refetchCategoriesIfRequested]);

	const isSignedIn = username !== null;
	return isSignedIn ? <SignedIn /> : <SignedOut />;
}
