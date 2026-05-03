import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute } from "@tanstack/react-router";
import { PlusIcon, XIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import z from "zod";

import { signInInstagram, editCategory, editUserCredentials } from "@/buf/raker/v1/raker-RakerServer_connectquery";
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
