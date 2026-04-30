import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { GalleryHorizontalIcon, Grid3x3Icon, TextAlignJustifyIcon, TrashIcon, ExternalLinkIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import z from "zod";

import { removeFiles, updateCategories } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import type { ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { FileDisplay, FilesCarousel, postTypeString } from "@/components/file-display";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useConfirmationDialog } from "@/hooks/use-confirmation-dialog";
import { useUser } from "@/hooks/user-provider";
import { cn } from "@/lib/utils";
import { HistoryPostCategoryForm } from "@/routes/history";

export function Result({ result }: { result: ScrapeResponse }) {
	const { username, categories: availableCategories } = useUser();
	const form = useForm({
		defaultValues: {
			categories: result.categories,
		},
		validators: {
			onSubmit: z.object({
				categories: z.array(z.string()).catch([]),
			}),
		},
		onSubmit: async ({ value: { categories } }) => {
			try {
				await updateCategoriesMutation.mutateAsync({
					type: currentResult.postType,
					owner: currentResult.postOwner,
					post: currentResult.post,
					categories,
				});

				setCurrentResult((previousResult) => {
					previousResult.categories = previousResult.categories.filter(
						(category) => !categories.includes(category),
					);
					return previousResult;
				});
				toast.success("Categories updated", {
					position: "top-center",
				});
			} catch (err) {
				toast.error((err as Error).message, {
					position: "top-center",
				});
			}
		},
	});
	const { confirm, DialogComponent } = useConfirmationDialog();
	const updateCategoriesMutation = useMutation(updateCategories);
	const removeFilesMutation = useMutation(removeFiles);
	const [currentResult, setCurrentResult] = useState(result);
	const [selection, setSelection] = useState<{ selectedFiles: string[]; anchorFile: string | null }>({
		selectedFiles: [],
		anchorFile: null,
	});
	const files = currentResult.files;

	useEffect(() => {
		setCurrentResult(result);
	}, [result]);

	useEffect(() => {
		form.setFieldValue("categories", currentResult.categories);
	}, [currentResult, form]);

	useEffect(() => {
		setSelection((current) => {
			const selectedFiles = files.filter((file) => current.selectedFiles.includes(file));
			const anchorFile =
				current.anchorFile !== null && files.includes(current.anchorFile)
					? current.anchorFile
					: (selectedFiles[0] ?? null);

			if (selectedFiles.length === current.selectedFiles.length && anchorFile === current.anchorFile) {
				return current;
			}

			return { selectedFiles, anchorFile };
		});
	}, [files]);

	useEffect(() => {
		const onKeyDown = (event: KeyboardEvent) => {
			if (event.key === "Escape") {
				setSelection({ selectedFiles: [], anchorFile: null });
			}
		};

		window.addEventListener("keydown", onKeyDown);
		return () => window.removeEventListener("keydown", onKeyDown);
	}, []);

	const isSelected = (file: string) => selection.selectedFiles.includes(file);

	const toggleSelection = (file: string) => {
		setSelection((current) => {
			const isCurrentlySelected = current.selectedFiles.includes(file);
			const selectedFiles = isCurrentlySelected
				? current.selectedFiles.filter((selectedFile) => selectedFile !== file)
				: [...current.selectedFiles, file];
			const anchorFile = isCurrentlySelected
				? current.anchorFile === file
					? (selectedFiles.at(-1) ?? null)
					: current.anchorFile
				: file;

			return { selectedFiles, anchorFile };
		});
	};

	const selectRange = (file: string) => {
		setSelection((current) => {
			if (current.anchorFile === null) {
				return { selectedFiles: [file], anchorFile: file };
			}

			const anchorIndex = files.indexOf(current.anchorFile);
			const fileIndex = files.indexOf(file);

			if (anchorIndex === -1 || fileIndex === -1) {
				return { selectedFiles: [file], anchorFile: file };
			}

			const start = Math.min(anchorIndex, fileIndex);
			const end = Math.max(anchorIndex, fileIndex);

			return { selectedFiles: files.slice(start, end + 1), anchorFile: current.anchorFile };
		});
	};

	const handleSelection = (file: string, event: { shiftKey: boolean; ctrlKey: boolean; metaKey: boolean }) => {
		if (event.shiftKey) {
			selectRange(file);
			return;
		}

		if (event.ctrlKey || event.metaKey) {
			toggleSelection(file);
			return;
		}

		toggleSelection(file);
	};

	const deleteFiles = async (paths: string[]) => {
		if (paths.length === 0) {
			return;
		}

		const confirmed = await confirm({
			title: "Delete Files",
			description: `Delete ${paths.length} file${paths.length === 1 ? "" : "s"}? This cannot be undone.`,
			confirmText: "Delete",
			cancelText: "Cancel",
			isDestructive: true,
		});

		if (!confirmed) {
			return;
		}

		try {
			const updatedResult = await removeFilesMutation.mutateAsync({
				type: currentResult.postType,
				owner: currentResult.postOwner,
				post: currentResult.post,
				paths,
			});

			setCurrentResult(updatedResult);
			setSelection({ selectedFiles: [], anchorFile: null });
		} catch (err) {
			toast.error((err as Error).message, {
				position: "top-center",
			});
		}
	};

	if (username === null) {
		return null;
	}

	return (
		<section className="my-3 flex w-full flex-col items-center gap-3">
			<Label>{timestampDate(currentResult.postDate!).toString()}</Label>
			<div className="w-full">
				<form.Field name="categories" mode="array">
					{(categoriesField) => (
						<div className="flex flex-col gap-3 rounded-lg border border-border/60 bg-background/70 p-3">
							<HistoryPostCategoryForm
								availableCategories={availableCategories}
								showExclusive={false}
								categoriesField={{
									name: categoriesField.name,
									value: categoriesField.state.value,
									onToggleCategory: (category, checked) => {
										if (checked) {
											if (!categoriesField.state.value.includes(category)) {
												categoriesField.pushValue(category);
											}
										} else {
											const index = categoriesField.state.value.indexOf(category);
											if (index > -1) {
												categoriesField.removeValue(index);
											}
										}
									},
								}}
							/>
							<div>
								<Button
									type="button"
									size="sm"
									disabled={updateCategoriesMutation.isPending}
									onClick={() => {
										form.handleSubmit();
									}}
								>
									Save Categories
								</Button>
							</div>
						</div>
					)}
				</form.Field>
			</div>
			<Tabs className="w-full">
				<div className="mx-auto flex w-full items-center gap-2 sm:w-1/2">
					<TabsList className="flex-1">
						<TabsTrigger value="list">
							<TextAlignJustifyIcon />
							List
						</TabsTrigger>
						<TabsTrigger value="grid">
							<Grid3x3Icon />
							Grid
						</TabsTrigger>
						<TabsTrigger value="carousel">
							<GalleryHorizontalIcon />
							Carousel
						</TabsTrigger>
					</TabsList>
					{selection.selectedFiles.length > 0 ? (
						<Button
							type="button"
							variant="destructive"
							size="sm"
							className="shrink-0"
							onClick={() => deleteFiles(selection.selectedFiles)}
						>
							Delete {selection.selectedFiles.length}
						</Button>
					) : null}
				</div>
				<TabsContent value="list">
					<Accordion multiple>
						{files.map((file) => {
							const selected = isSelected(file);

							return (
								<AccordionItem
									aria-pressed={selected}
									className={cn(
										"flex w-full flex-col rounded-lg border border-transparent px-2 py-1 transition",
										selected && "border-primary/60 bg-primary/10 ring-2 ring-primary/30",
									)}
									key={`accordion-file-${currentResult.postType}-${currentResult.postOwner}-${currentResult.post}-${file}`}
									value={file}
								>
									<div className="flex w-full items-center gap-2">
										<Checkbox
											checked={selected}
											aria-label={selected ? `Deselect ${file}` : `Select ${file}`}
											onClick={(event) => {
												event.preventDefault();
												event.stopPropagation();
												handleSelection(file, event);
											}}
										/>
										<Button
											type="button"
											variant="outline"
											size="sm"
											className="hover:bg-blue/20 hover:text-blue shrink-0 px-2"
											render={
												<a
													href={`/api/storage/${username}/${postTypeString(result.postType)}/${result.postOwner}/${file}`}
													target="_blank"
													rel="noopener noreferrer"
													aria-label={`Open ${file} in new tab`}
												/>
											}
										>
											<ExternalLinkIcon className="h-4 w-4" />
										</Button>
										<Button
											type="button"
											variant="outline"
											size="sm"
											className="shrink-0 px-2 hover:bg-destructive/20 hover:text-destructive"
											disabled={removeFilesMutation.isPending}
											onClick={() => deleteFiles([file])}
											aria-label={`Delete ${file}`}
										>
											<TrashIcon className="h-4 w-4" />
										</Button>
										<AccordionTrigger className="flex-1 gap-2 text-left">
											<Label className="w-full wrap-anywhere whitespace-normal">{file}</Label>
										</AccordionTrigger>
									</div>
									<AccordionContent className="max-w-[50vw]">
										<FileDisplay file={file} post={currentResult} username={username} />
									</AccordionContent>
								</AccordionItem>
							);
						})}
					</Accordion>
				</TabsContent>
				<TabsContent value="grid" className="grid grid-cols-2 gap-3 lg:grid-cols-3">
					{files.map((file) => {
						const selected = isSelected(file);

						return (
							<div
								aria-pressed={selected}
								className={cn(
									"relative rounded-xl border border-transparent p-1 transition",
									selected && "border-primary/60 bg-primary/10 ring-2 ring-primary/30",
								)}
								key={`grid-file-${currentResult.postType}-${currentResult.postOwner}-${currentResult.post}-${file}`}
							>
								<Checkbox
									checked={selected}
									aria-label={selected ? `Deselect ${file}` : `Select ${file}`}
									className="absolute top-2 right-2 z-10 border-foreground/60 bg-background/95 shadow-sm dark:border-foreground/70 dark:bg-background/85"
									onClick={(event) => {
										event.preventDefault();
										event.stopPropagation();
										handleSelection(file, event);
									}}
								/>
								<FileDisplay file={file} post={currentResult} username={username} />
							</div>
						);
					})}
				</TabsContent>
				<TabsContent
					value="carousel"
					className="mt-2 w-full [&_img]:max-h-[50vh] [&_img]:w-auto [&_video]:max-h-[50vh] [&_video]:w-auto"
				>
					<FilesCarousel post={currentResult} username={username} />
				</TabsContent>
			</Tabs>
			<DialogComponent />
		</section>
	);
}
