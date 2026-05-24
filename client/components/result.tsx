import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { Link } from "@tanstack/react-router";
import {
	GalleryHorizontalIcon,
	Grid3x3Icon,
	TextAlignJustifyIcon,
	TrashIcon,
	ExternalLinkIcon,
	CopyIcon,
	CropIcon,
	ImageIcon,
} from "lucide-react";
import { useEffect, useState, type Dispatch, type ReactNode, type SetStateAction } from "react";
import { toast } from "sonner";
import z from "zod";

import { removeFiles, updateCategories } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { FileDisplay, FilesCarousel, FileSheet, postTypeString } from "@/components/file-display";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
	ContextMenu,
	ContextMenuContent,
	ContextMenuGroup,
	ContextMenuItem,
	ContextMenuTrigger,
} from "@/components/ui/context-menu";
import { Label } from "@/components/ui/label";
import { Progress } from "@/components/ui/progress";
import { GoogleMapsLink } from "@/components/ui/svgs/google-maps";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { SnapchatIcon } from "@/components/ui/svgs/snapchat";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useConfirmationDialog } from "@/hooks/use-confirmation-dialog";
import { useUser } from "@/hooks/user-provider";
import { cn, writeClipboard, defaultPostTypes, inPWA, uniqueArraysEqualAsSets } from "@/lib/utils";
import { HistoryPostCategoryForm } from "@/routes/history";

export function PlatformIcon({ type }: { type: PostType | -1 }) {
	switch (type) {
		case PostType.Instagram:
		case PostType.Highlight:
		case PostType.Story:
			return <InstagramIcon className="w-4" />;
		case PostType.TikTok:
			return <TikTokIcon className="w-4" />;
		case PostType.Snapchat:
			return <SnapchatIcon className="w-4" />;
		case PostType.VSCO:
			return <VSCOIcon className="w-4" />;
		default:
			return null;
	}
}

export function ResultLink({ result, children }: { result: ScrapeResponse; children: ReactNode }) {
	const target = inPWA() ? undefined : "_blank";

	switch (result.postType) {
		case PostType.Instagram:
			return (
				<Link to="/instagram" search={{ post: result.post, incognito: result.incognito }} target={target}>
					{children}
				</Link>
			);
		case PostType.Highlight:
			return (
				<Link to="/highlight" search={{ highlight: result.post }} target={target}>
					{children}
				</Link>
			);
		case PostType.Story:
			return (
				<Link to="/story" search={{ owner: result.post }} target={target}>
					{children}
				</Link>
			);
		case PostType.TikTok:
			return (
				<Link
					to="/tiktok"
					search={{ owner: result.postOwner, post: result.post, incognito: result.incognito }}
					target={target}
				>
					{children}
				</Link>
			);
		case PostType.Snapchat:
			return (
				<Link to="/snapchat" search={{ owner: result.postOwner, highlight: result.post }} target={target}>
					{children}
				</Link>
			);
		case PostType.VSCO:
			return (
				<Link to="/vsco" search={{ owner: result.postOwner, post: result.post }} target={target}>
					{children}
				</Link>
			);
	}
}

export function PostTypeIconLabel({ type }: { type: PostType }) {
	switch (type) {
		case PostType.Instagram:
			return (
				<span className="inline-flex items-center gap-1 align-middle leading-none whitespace-nowrap">
					<InstagramIcon className="w-4" />
					Post
				</span>
			);
		case PostType.Highlight:
			return (
				<span className="inline-flex items-center gap-1 align-middle leading-none whitespace-nowrap">
					<InstagramIcon className="w-4" />
					Highlight
				</span>
			);
		case PostType.Story:
			return (
				<span className="inline-flex items-center gap-1 align-middle leading-none whitespace-nowrap">
					<InstagramIcon className="w-4" />
					Story
				</span>
			);
		case PostType.TikTok:
			return (
				<span className="inline-flex items-center gap-1 align-middle leading-none whitespace-nowrap">
					<TikTokIcon className="w-4" />
					Post
				</span>
			);
		case PostType.Snapchat:
			return (
				<span className="inline-flex items-center gap-1 align-middle leading-none whitespace-nowrap">
					<SnapchatIcon className="w-4" />
					Highlight
				</span>
			);
		case PostType.VSCO:
			return (
				<span className="inline-flex items-center gap-1 align-middle leading-none whitespace-nowrap">
					<VSCOIcon className="w-4" />
					Post
				</span>
			);
	}
}

export function ResultHeader({
	result,
	categories,
	exclusive,
	showPost,
}: {
	result: ScrapeResponse;
	categories: string[];
	exclusive: boolean;
	showPost?: boolean;
}) {
	const target = inPWA() ? undefined : "_blank";

	return (
		<span className="inline-flex max-w-full flex-wrap items-center gap-x-1 gap-y-1 leading-none">
			<Badge variant="secondary">
				<Link
					to="/history"
					search={{
						categories,
						exclusive,
						page: 1n,
						owners: [],
						types: [result.postType],
					}}
					target={target}
				>
					<PostTypeIconLabel type={result.postType} />
				</Link>
			</Badge>
			<span>/</span>
			<ContextMenu>
				<ContextMenuTrigger>
					<Badge variant="secondary">
						<code className="align-middle leading-none">{result.postOwner}</code>
					</Badge>
				</ContextMenuTrigger>
				<ContextMenuContent>
					<ContextMenuGroup>
						<Link
							to="/history"
							search={{
								categories,
								exclusive,
								page: 1n,
								owners: [{ owner: result.postOwner, type: -1 }],
								types: defaultPostTypes,
							}}
							target={target}
						>
							<ContextMenuItem>
								<ExternalLinkIcon /> History Results
							</ContextMenuItem>
						</Link>
						<a
							target="_blank"
							rel="noopener noreferrer"
							href={(() => {
								switch (result.postType) {
									case PostType.Instagram:
									case PostType.Highlight:
									case PostType.Story:
										return `https://www.instagram.com/${result.postOwner}`;
									case PostType.TikTok:
										return `https://www.tiktok.com/@${result.postOwner}`;
									case PostType.Snapchat:
										return `https://www.snapchat.com/@${result.postOwner}`;
									case PostType.VSCO:
										return `https://vsco.co/${result.postOwner}/gallery`;
								}
							})()}
						>
							<ContextMenuItem>
								<PlatformIcon type={result.postType} /> Open Profile
							</ContextMenuItem>
						</a>
						<ContextMenuItem onClick={() => writeClipboard(result.postOwner)}>
							<CopyIcon /> Copy
						</ContextMenuItem>
					</ContextMenuGroup>
				</ContextMenuContent>
			</ContextMenu>
			<span>/</span>
			{showPost ? (
				<ContextMenu>
					<ContextMenuTrigger>
						<Badge variant="secondary">
							<code className="align-middle leading-none">{result.post}</code>
						</Badge>
					</ContextMenuTrigger>
					<ContextMenuContent>
						<ContextMenuGroup>
							<ResultLink result={result}>
								<ContextMenuItem>
									<ExternalLinkIcon /> History Result
								</ContextMenuItem>
							</ResultLink>
							{![PostType.Highlight, PostType.Story].includes(result.postType) && (
								<a
									target="_blank"
									rel="noopener noreferrer"
									href={(() => {
										switch (result.postType) {
											case PostType.Instagram:
											case PostType.Highlight:
											case PostType.Story:
												return `https://www.instagram.com/p/${result.post}`;
											case PostType.TikTok:
												return `https://www.tiktok.com/@${result.postOwner}/video/${result.post}`;
											case PostType.Snapchat:
												return `https://www.snapchat.com/@${result.postOwner}/highlight/${result.post}`;
											case PostType.VSCO:
												return `https://vsco.co/${result.postOwner}/media/${result.post}`;
										}
									})()}
								>
									<ContextMenuItem>
										<PlatformIcon type={result.postType} /> Open Original Post
									</ContextMenuItem>
								</a>
							)}
						</ContextMenuGroup>
					</ContextMenuContent>
				</ContextMenu>
			) : (
				<Badge variant="secondary">
					<code className="align-middle leading-none">{result.post}</code>
				</Badge>
			)}
		</span>
	);
}

export function Result({
	result,
	setResult,
}: {
	result: ScrapeResponse;
	setResult: Dispatch<SetStateAction<ScrapeResponse | null>>;
}) {
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
					type: result.postType,
					owner: result.postOwner,
					post: result.post,
					categories,
				});

				setResult((previousResult) => {
					if (previousResult === null) {
						return previousResult;
					}

					return {
						...previousResult,
						categories,
					};
				});
				toast.success("Updated", {
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
	const [selection, setSelection] = useState<{ selectedFiles: string[]; anchorFile: string | null }>({
		selectedFiles: [],
		anchorFile: null,
	});
	const files = result.files;

	useEffect(() => {
		form.setFieldValue("categories", result.categories);
	}, [form, result]);

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
				type: result.postType,
				owner: result.postOwner,
				post: result.post,
				paths,
			});

			setResult(updatedResult.files.length === 0 ? null : updatedResult);
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
			{(updateCategoriesMutation.isPending || removeFilesMutation.isPending) && (
				<div className="w-full">
					<Progress value={null} className="pb-2" />
				</div>
			)}
			<div className="max-w-full">
				<ResultHeader categories={availableCategories} exclusive={false} result={result} />
			</div>
			<Label>{timestampDate(result.postDate!).toString()}</Label>
			<div className="w-full">
				<form.Field name="categories" mode="array">
					{(categoriesField) => {
						const hasUnsavedCategories = !uniqueArraysEqualAsSets(
							categoriesField.state.value,
							result.categories,
						);
						return (
							<div className="flex flex-col gap-3 rounded-lg border border-border/60 bg-background/70 p-3">
								<HistoryPostCategoryForm
									availableCategories={availableCategories}
									showExclusive={false}
									legendBadge={
										hasUnsavedCategories ? (
											<Badge className="mr-1 h-2 w-2 rounded-full p-0" />
										) : null
									}
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
										disabled={updateCategoriesMutation.isPending || !hasUnsavedCategories}
										onClick={() => {
											form.handleSubmit();
										}}
									>
										Save Categories
									</Button>
								</div>
							</div>
						);
					}}
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
						{result.files.length > 1 && (
							<TabsTrigger value="carousel">
								<GalleryHorizontalIcon />
								Carousel
							</TabsTrigger>
						)}
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
									key={`accordion-file-${result.postType}-${result.postOwner}-${result.post}-${file}`}
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
											nativeButton={false}
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
										{/\.(jpe?g)$/.test(file) && (
											<FileSheet
												file={file}
												post={result}
												username={username}
												trigger={
													<Button
														variant="outline"
														size="sm"
														className="dark:bg-secondary dark:hover:bg-secondary/80"
													>
														<ImageIcon />
														<CropIcon />
													</Button>
												}
											/>
										)}
										{result.postType === PostType.VSCO && result.coordinates ? (
											<GoogleMapsLink coordinates={result.coordinates} size="sm" />
										) : null}
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
									<AccordionContent className="sm:max-w-full md:max-w-[25vw]">
										<FileDisplay file={file} post={result} username={username} />
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
								key={`grid-file-${result.postType}-${result.postOwner}-${result.post}-${file}`}
							>
								<Checkbox
									checked={selected}
									aria-label={selected ? `Deselect ${file}` : `Select ${file}`}
									className="absolute top-3 right-3 z-10 border-foreground/60 bg-background/95 shadow-sm dark:border-foreground/70 dark:bg-background/85"
									onClick={(event) => {
										event.preventDefault();
										event.stopPropagation();
										handleSelection(file, event);
									}}
								/>
								<FileDisplay file={file} post={result} username={username} withCrop withCoordinates />
							</div>
						);
					})}
				</TabsContent>
				{result.files.length > 1 && (
					<TabsContent
						value="carousel"
						className="mt-2 w-full [&_img]:max-h-[50vh] [&_img]:w-auto [&_video]:max-h-[50vh] [&_video]:w-auto"
					>
						<FilesCarousel post={result} username={username} />
					</TabsContent>
				)}
			</Tabs>
			<DialogComponent />
		</section>
	);
}
