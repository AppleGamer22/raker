import { useMutation } from "@connectrpc/connect-query";
import { CheckIcon, CropIcon, ImageIcon, UndoIcon } from "lucide-react";
import {
	useCallback,
	useEffect,
	useRef,
	useState,
	type CSSProperties,
	type ReactElement,
	type SyntheticEvent,
	type PointerEvent as ReactPointerEvent,
} from "react";

import "react-resizable/css/styles.css";

import { ResizableBox, type ResizeCallbackData } from "react-resizable";
import { toast } from "sonner";

import { cropFile } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { Button } from "@/components/ui/button";
import { ButtonGroup } from "@/components/ui/button-group";
import {
	Carousel,
	CarouselContent,
	CarouselItem,
	CarouselNext,
	CarouselPrevious,
	type CarouselApi,
} from "@/components/ui/carousel";
import { Progress } from "@/components/ui/progress";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { GoogleMapsLink } from "@/components/ui/svgs/google-maps";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";

export function postTypeString(type: PostType): string {
	switch (type) {
		case PostType.Instagram:
			return "instagram";
		case PostType.Highlight:
			return "highlight";
		case PostType.Story:
			return "story";
		case PostType.TikTok:
			return "tiktok";
		case PostType.Snapchat:
			return "snapchat";
		case PostType.VSCO:
			return "vsco";
	}
}

export function FilesCarousel({
	username,
	post: { postType, postOwner, coordinates, post, files },
}: {
	username: string;
	post: ScrapeResponse;
}) {
	const [api, setApi] = useState<CarouselApi>();
	const [selectedIndex, setSelectedIndex] = useState(0);
	const prevHeightsRef = useRef<Map<number, number>>(new Map());

	const syncSelectedSlideHeight = useCallback((emblaApi: CarouselApi) => {
		if (!emblaApi) {
			return;
		}

		const selectedIndex = emblaApi.selectedScrollSnap();
		const selectedSlide = emblaApi.slideNodes()[selectedIndex];
		if (!selectedSlide) {
			return;
		}

		const h = selectedSlide.offsetHeight;
		const prev = prevHeightsRef.current.get(selectedIndex) || 0;

		if (h > prev) {
			prevHeightsRef.current.set(selectedIndex, h);
			emblaApi.containerNode().style.height = `${h}px`;
		} else if (prev > 0) {
			// Do not shrink below the previously observed max height for this slide.
			emblaApi.containerNode().style.height = `${prev}px`;
		} else {
			emblaApi.containerNode().style.height = `${h}px`;
		}
	}, []);

	useEffect(() => {
		if (!api) {
			return;
		}

		const onSelect = () => {
			setSelectedIndex(api.selectedScrollSnap());
			syncSelectedSlideHeight(api);
		};

		const onReInit = () => {
			setSelectedIndex(api.selectedScrollSnap());
			syncSelectedSlideHeight(api);
		};

		onSelect();
		api.on("reInit", onReInit);
		api.on("select", onSelect);
		api.on("settle", onSelect);

		return () => {
			api.off("reInit", onReInit);
			api.off("select", onSelect);
			api.off("settle", onSelect);
		};
	}, [api, syncSelectedSlideHeight]);

	return files.length > 1 ? (
		<Carousel opts={{ loop: true }} setApi={setApi}>
			<CarouselContent className="items-center">
				{files.map((file, i) => (
					<CarouselItem
						key={`file-${postType}-${postOwner}-${post}-${i}`}
						className="flex items-center justify-center self-center"
					>
						<FileDisplay
							username={username}
							post={{ postType, postOwner, coordinates } as ScrapeResponse}
							file={file}
							onMediaLoad={() => {
								api?.reInit();
								if (api) {
									requestAnimationFrame(() => syncSelectedSlideHeight(api));
								}
							}}
							withCrop
							withCoordinates
						/>
					</CarouselItem>
				))}
			</CarouselContent>
			<CarouselPrevious className="top-1/2 left-2 z-10 -translate-y-1/2" />
			<CarouselNext className="top-1/2 right-2 z-10 -translate-y-1/2" />
			<div className="flex items-center justify-center gap-2 pt-2" role="group" aria-label="Slide navigation">
				{files.map((_, index) => (
					<button
						type="button"
						key={`dot-${postType}-${postOwner}-${post}-${index}`}
						onClick={() => api?.scrollTo(index)}
						className={`h-2 rounded-full transition-all focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none ${
							index === selectedIndex
								? "w-6 bg-primary"
								: "w-2 bg-muted-foreground/40 hover:bg-muted-foreground/70"
						}`}
						aria-label={`Go to slide ${index + 1}`}
						aria-current={index === selectedIndex ? "true" : undefined}
					/>
				))}
			</div>
		</Carousel>
	) : (
		<FileDisplay
			username={username}
			post={{ postType, postOwner, coordinates } as ScrapeResponse}
			file={files[0]}
			withCrop
			withCoordinates
		/>
	);
}

export function FileDisplay({
	username,
	file,
	post: { postType, postOwner, coordinates },
	onMediaLoad,
	withCrop,
	withCoordinates,
	className = "h-auto w-full rounded-xl",
	cacheBuster,
}: {
	username: string;
	file: string;
	post: ScrapeResponse;
	onMediaLoad?: () => void;
	withCrop?: boolean;
	withCoordinates?: boolean;
	className?: string;
	cacheBuster?: number | string;
}) {
	const [cacheBusterState, setCacheBusterState] = useState<number | string | undefined>(cacheBuster ?? undefined);

	useEffect(() => {
		setCacheBusterState(cacheBuster ?? undefined);
	}, [cacheBuster]);

	useEffect(() => {
		const handler = (ev: Event) => {
			try {
				const custom = ev as CustomEvent<any>;
				const d = custom.detail;
				if (
					d &&
					d.username === username &&
					d.postType === postTypeString(postType) &&
					d.postOwner === postOwner &&
					d.file === file
				) {
					setCacheBusterState(d.cacheBuster);
				}
			} catch {}
		};
		window.addEventListener("fileCropped", handler as EventListener);
		return () => window.removeEventListener("fileCropped", handler as EventListener);
	}, [username, postType, postOwner, file]);

	const url =
		`/api/storage/${username}/${postTypeString(postType)}/${postOwner}/${file}` +
		(cacheBusterState ? `?v=${cacheBusterState}` : "");
	if (/\.(jpe?g)|(webp)|(heic)$/.test(file)) {
		const imgResult = <img src={url} onLoad={onMediaLoad} loading="lazy" className={className} />;
		return withCrop || withCoordinates ? (
			<div className="relative inline-block w-full rounded-xl">
				{imgResult}
				<div className="absolute top-2 left-2 z-10 flex flex-row gap-2">
					{withCrop && /\.(jpe?g)$/.test(file) && (
						<FileSheet
							file={file}
							post={{ postType, postOwner, coordinates } as ScrapeResponse}
							username={username}
							trigger={
								<Button variant="outline" className="dark:bg-secondary dark:hover:bg-secondary/80">
									<ImageIcon />
									<CropIcon />
								</Button>
							}
						/>
					)}
					{withCoordinates && postType === PostType.VSCO && coordinates && (
						<GoogleMapsLink coordinates={coordinates} size="icon" />
					)}
				</div>
			</div>
		) : (
			imgResult
		);
	} else if (/\.(mp4)|(webm)$/.test(file)) {
		return (
			<video
				src={url}
				onLoadedMetadata={onMediaLoad}
				preload="metadata"
				className={className}
				loop
				controls
				muted
			/>
		);
	} else {
		return <a href={url}>{url}</a>;
	}
}

export type CropRect = {
	x1: number;
	y1: number;
	x2: number;
	y2: number;
};

export type CropBox = {
	x: number;
	y: number;
	width: number;
	height: number;
};

const CROP_HANDLE_SIZE = 10;
const MIN_CROP_SIZE = 40;
const FULL_IMAGE_CROP_EPSILON = 1;

function clamp(value: number, min: number, max: number) {
	return Math.min(Math.max(value, min), max);
}

function clampRect(rect: CropRect, maxWidth: number, maxHeight: number): CropRect {
	const left = Math.min(rect.x1, rect.x2);
	const right = Math.max(rect.x1, rect.x2);
	const top = Math.min(rect.y1, rect.y2);
	const bottom = Math.max(rect.y1, rect.y2);
	const safeX1 = clamp(left, 0, Math.max(0, maxWidth - MIN_CROP_SIZE));
	const safeY1 = clamp(top, 0, Math.max(0, maxHeight - MIN_CROP_SIZE));
	const safeX2 = clamp(right, safeX1 + MIN_CROP_SIZE, maxWidth);
	const safeY2 = clamp(bottom, safeY1 + MIN_CROP_SIZE, maxHeight);
	return {
		x1: safeX1,
		y1: safeY1,
		x2: safeX2,
		y2: safeY2,
	};
}

function clampRectPreserveSize(rect: CropRect, maxWidth: number, maxHeight: number): CropRect {
	const left = Math.min(rect.x1, rect.x2);
	const right = Math.max(rect.x1, rect.x2);
	const top = Math.min(rect.y1, rect.y2);
	const bottom = Math.max(rect.y1, rect.y2);
	const width = Math.max(MIN_CROP_SIZE, right - left);
	const height = Math.max(MIN_CROP_SIZE, bottom - top);
	const safeX1 = clamp(left, 0, Math.max(0, maxWidth - width));
	const safeY1 = clamp(top, 0, Math.max(0, maxHeight - height));
	return {
		x1: safeX1,
		y1: safeY1,
		x2: safeX1 + width,
		y2: safeY1 + height,
	};
}

function displayToNatural(rect: CropRect, scaleX: number, scaleY: number): CropRect {
	return {
		x1: Math.round(rect.x1 * scaleX),
		y1: Math.round(rect.y1 * scaleY),
		x2: Math.round(rect.x2 * scaleX),
		y2: Math.round(rect.y2 * scaleY),
	};
}

function naturalToDisplay(rect: CropRect, scaleX: number, scaleY: number): CropRect {
	return {
		x1: rect.x1 / scaleX,
		y1: rect.y1 / scaleY,
		x2: rect.x2 / scaleX,
		y2: rect.y2 / scaleY,
	};
}

export function cropRectToBox(rect: CropRect): CropBox {
	const x1 = Math.min(rect.x1, rect.x2);
	const y1 = Math.min(rect.y1, rect.y2);
	const x2 = Math.max(rect.x1, rect.x2);
	const y2 = Math.max(rect.y1, rect.y2);
	return {
		x: x1,
		y: y1,
		width: x2 - x1,
		height: y2 - y1,
	};
}

function boxToRect(box: { x: number; y: number; width: number; height: number }): CropRect {
	return {
		x1: box.x,
		y1: box.y,
		x2: box.x + box.width,
		y2: box.y + box.height,
	};
}

function isFullImageCrop(rect: CropRect, naturalSize: { width: number; height: number }): boolean {
	return (
		rect.x1 <= FULL_IMAGE_CROP_EPSILON &&
		rect.y1 <= FULL_IMAGE_CROP_EPSILON &&
		rect.x2 >= naturalSize.width - FULL_IMAGE_CROP_EPSILON &&
		rect.y2 >= naturalSize.height - FULL_IMAGE_CROP_EPSILON
	);
}

function handleStyle(handle: string): CSSProperties {
	const base: CSSProperties = {
		position: "absolute",
		width: CROP_HANDLE_SIZE,
		height: CROP_HANDLE_SIZE,
	};

	switch (handle) {
		case "se":
			return { ...base, right: 0, bottom: 0 };
		case "sw":
			return { ...base, left: 0, bottom: 0 };
		case "ne":
			return { ...base, right: 0, top: 0 };
		case "nw":
			return { ...base, left: 0, top: 0 };
		case "n":
			return { ...base, top: 0, left: "50%", transform: "translate(-50%, 0)" };
		case "s":
			return { ...base, bottom: 0, left: "50%", transform: "translate(-50%, 0)" };
		case "e":
			return { ...base, right: 0, top: "50%", transform: "translate(0, -50%)" };
		case "w":
			return { ...base, left: 0, top: "50%", transform: "translate(0, -50%)" };
		default:
			return base;
	}
}

function CropPreview({
	username,
	file,
	post: { postType, postOwner },
	onCropChange,
	className = "h-full w-auto rounded-xl",
	cacheBuster,
	resetSignal,
}: {
	username: string;
	file: string;
	post: ScrapeResponse;
	onCropChange?: (rect: CropRect | null, isFullImageCrop: boolean) => void;
	className?: string;
	cacheBuster?: number | string;
	resetSignal?: number;
}) {
	const imageRef = useRef<HTMLImageElement | null>(null);
	const dragStateRef = useRef<{
		startX: number;
		startY: number;
		originX: number;
		originY: number;
		width: number;
		height: number;
	} | null>(null);
	const resizeHandleRef = useRef<string | null>(null);
	const latestCropNaturalRef = useRef<CropRect | null>(null);
	const [activeHandle, setActiveHandle] = useState<string | null>(null);
	const [naturalSize, setNaturalSize] = useState({ width: 0, height: 0 });
	const [displaySize, setDisplaySize] = useState({ width: 0, height: 0 });
	const [cropNatural, setCropNatural] = useState<CropRect | null>(null);

	const updateDisplaySize = useCallback(() => {
		const img = imageRef.current;
		if (!img) {
			return;
		}
		const rect = img.getBoundingClientRect();
		const width = Math.round(rect.width);
		const height = Math.round(rect.height);
		setDisplaySize((prev) => (prev.width === width && prev.height === height ? prev : { width, height }));
	}, []);

	useEffect(() => {
		setCropNatural(null);
		setNaturalSize({ width: 0, height: 0 });
		setDisplaySize({ width: 0, height: 0 });
		latestCropNaturalRef.current = null;
	}, [file]);

	useEffect(() => {
		const img = imageRef.current;
		if (!img || typeof ResizeObserver === "undefined") {
			return;
		}
		const observer = new ResizeObserver(() => updateDisplaySize());
		observer.observe(img);
		return () => observer.disconnect();
	}, [updateDisplaySize, file]);

	const handleImageLoad = useCallback(() => {
		const img = imageRef.current;
		if (!img) {
			return;
		}
		setNaturalSize({ width: img.naturalWidth, height: img.naturalHeight });
		updateDisplaySize();
	}, [updateDisplaySize]);

	const scaleX = naturalSize.width > 0 && displaySize.width > 0 ? naturalSize.width / displaySize.width : 1;
	const scaleY = naturalSize.height > 0 && displaySize.height > 0 ? naturalSize.height / displaySize.height : 1;
	const getFullDisplayRect = useCallback(() => {
		return clampRect(
			{
				x1: 0,
				y1: 0,
				x2: displaySize.width,
				y2: displaySize.height,
			},
			displaySize.width,
			displaySize.height,
		);
	}, [displaySize.height, displaySize.width]);

	useEffect(() => {
		if (
			cropNatural ||
			displaySize.width === 0 ||
			displaySize.height === 0 ||
			naturalSize.width === 0 ||
			naturalSize.height === 0
		) {
			return;
		}
		const initialDisplay = getFullDisplayRect();
		const initialNatural = displayToNatural(initialDisplay, scaleX, scaleY);
		latestCropNaturalRef.current = initialNatural;
		setCropNatural(initialNatural);
	}, [cropNatural, displaySize, naturalSize, scaleX, scaleY, getFullDisplayRect]);

	const cropDisplay = cropNatural
		? clampRect(naturalToDisplay(cropNatural, scaleX, scaleY), displaySize.width, displaySize.height)
		: null;
	const cropDisplayBox = cropDisplay ? cropRectToBox(cropDisplay) : null;
	const edgeContact =
		!!cropDisplayBox && displaySize.width > 0 && displaySize.height > 0
			? (() => {
					const threshold = 1;
					const left = cropDisplayBox.x <= threshold;
					const right = cropDisplayBox.x + cropDisplayBox.width >= displaySize.width - threshold;
					const top = cropDisplayBox.y <= threshold;
					const bottom = cropDisplayBox.y + cropDisplayBox.height >= displaySize.height - threshold;
					return { left, right, top, bottom };
				})()
			: { left: false, right: false, top: false, bottom: false };
	const maxConstraints: [number, number] = cropDisplay
		? (() => {
				const handle = activeHandle ?? resizeHandleRef.current ?? "se";
				const maxWidth = handle.includes("w")
					? (cropDisplayBox?.x ?? 0) + (cropDisplayBox?.width ?? 0)
					: displaySize.width - (cropDisplayBox?.x ?? 0);
				const maxHeight = handle.includes("n")
					? (cropDisplayBox?.y ?? 0) + (cropDisplayBox?.height ?? 0)
					: displaySize.height - (cropDisplayBox?.y ?? 0);
				return [Math.max(MIN_CROP_SIZE, maxWidth), Math.max(MIN_CROP_SIZE, maxHeight)];
			})()
		: [MIN_CROP_SIZE, MIN_CROP_SIZE];

	const updateFromDisplay = useCallback(
		(displayRect: CropRect, options?: { preserveSize?: boolean }) => {
			if (displaySize.width === 0 || displaySize.height === 0) {
				return;
			}
			const clamped = options?.preserveSize
				? clampRectPreserveSize(displayRect, displaySize.width, displaySize.height)
				: clampRect(displayRect, displaySize.width, displaySize.height);
			const nextNatural = displayToNatural(clamped, scaleX, scaleY);
			latestCropNaturalRef.current = nextNatural;
			setCropNatural(nextNatural);
		},
		[displaySize, scaleX, scaleY],
	);

	const emitCropChange = useCallback(
		(rect: CropRect | null) => {
			if (!onCropChange) {
				return;
			}
			if (!rect || naturalSize.width === 0 || naturalSize.height === 0) {
				onCropChange(null, true);
				return;
			}
			onCropChange(rect, isFullImageCrop(rect, naturalSize));
		},
		[onCropChange, naturalSize],
	);

	useEffect(() => {
		if (
			resetSignal === undefined ||
			displaySize.width === 0 ||
			displaySize.height === 0 ||
			naturalSize.width === 0 ||
			naturalSize.height === 0
		) {
			return;
		}
		const fullDisplay = getFullDisplayRect();
		const fullNatural = displayToNatural(fullDisplay, scaleX, scaleY);
		latestCropNaturalRef.current = fullNatural;
		setCropNatural(fullNatural);
		emitCropChange(fullNatural);
	}, [resetSignal, displaySize, naturalSize, scaleX, scaleY, emitCropChange, getFullDisplayRect]);

	const handleResize = useCallback(
		(_event: SyntheticEvent, data: ResizeCallbackData) => {
			if (!cropDisplay || !cropDisplayBox) {
				return;
			}
			const handle = data.handle ?? resizeHandleRef.current ?? "se";
			const nextWidth = data.size.width;
			const nextHeight = data.size.height;
			const shiftX = cropDisplayBox.width - nextWidth;
			const shiftY = cropDisplayBox.height - nextHeight;
			const nextX = handle.includes("w") ? cropDisplayBox.x + shiftX : cropDisplayBox.x;
			const nextY = handle.includes("n") ? cropDisplayBox.y + shiftY : cropDisplayBox.y;
			updateFromDisplay(boxToRect({ x: nextX, y: nextY, width: nextWidth, height: nextHeight }));
		},
		[cropDisplay, cropDisplayBox, updateFromDisplay],
	);

	const handleResizeStart = useCallback((_event: SyntheticEvent, data: ResizeCallbackData) => {
		const handle = data.handle ?? null;
		resizeHandleRef.current = handle;
		setActiveHandle(handle);
	}, []);

	const handleResizeStop = useCallback(() => {
		resizeHandleRef.current = null;
		setActiveHandle(null);
		emitCropChange(latestCropNaturalRef.current);
	}, [emitCropChange]);

	const handleDragStart = useCallback(
		(event: ReactPointerEvent<HTMLDivElement>) => {
			if (!cropDisplay || !cropDisplayBox || event.button !== 0) {
				return;
			}
			event.preventDefault();

			const target = event.currentTarget;
			if (target.setPointerCapture) {
				target.setPointerCapture(event.pointerId);
			}
			dragStateRef.current = {
				startX: event.clientX,
				startY: event.clientY,
				originX: cropDisplayBox.x,
				originY: cropDisplayBox.y,
				width: cropDisplayBox.width,
				height: cropDisplayBox.height,
			};

			const handleMove = (moveEvent: PointerEvent) => {
				const dragState = dragStateRef.current;
				if (!dragState) {
					return;
				}
				const dx = moveEvent.clientX - dragState.startX;
				const dy = moveEvent.clientY - dragState.startY;
				updateFromDisplay(
					{
						x1: dragState.originX + dx,
						y1: dragState.originY + dy,
						x2: dragState.originX + dx + dragState.width,
						y2: dragState.originY + dy + dragState.height,
					},
					{ preserveSize: true },
				);
			};

			const handleUp = () => {
				dragStateRef.current = null;
				target.removeEventListener("pointermove", handleMove);
				target.removeEventListener("pointerup", handleUp);
				target.removeEventListener("pointercancel", handleUp);
				if (target.releasePointerCapture) {
					target.releasePointerCapture(event.pointerId);
				}
				emitCropChange(latestCropNaturalRef.current);
			};

			target.addEventListener("pointermove", handleMove);
			target.addEventListener("pointerup", handleUp);
			target.addEventListener("pointercancel", handleUp);
		},
		[cropDisplay, cropDisplayBox, updateFromDisplay, emitCropChange],
	);

	const url =
		`/api/storage/${username}/${postTypeString(postType)}/${postOwner}/${file}` +
		(cacheBuster ? `?v=${cacheBuster}` : "");

	return (
		<div className={cn("relative inline-block", className)}>
			<img ref={imageRef} src={url} onLoad={handleImageLoad} className={className} />
			{cropDisplayBox && (
				<ResizableBox
					width={cropDisplayBox.width}
					height={cropDisplayBox.height}
					className="crop-box absolute z-10 rounded-md border-2 border-primary/80"
					style={{
						left: cropDisplayBox.x,
						top: cropDisplayBox.y,
						position: "absolute",
						borderTopColor: edgeContact.top ? "hsl(var(--secondary))" : undefined,
						borderBottomColor: edgeContact.bottom ? "hsl(var(--secondary))" : undefined,
						borderLeftColor: edgeContact.left ? "hsl(var(--secondary))" : undefined,
						borderRightColor: edgeContact.right ? "hsl(var(--secondary))" : undefined,
					}}
					handleSize={[CROP_HANDLE_SIZE, CROP_HANDLE_SIZE]}
					minConstraints={[MIN_CROP_SIZE, MIN_CROP_SIZE]}
					maxConstraints={maxConstraints}
					handle={(handle, ref) => (
						<span
							ref={ref}
							className={`react-resizable-handle react-resizable-handle-${handle} crop-handle`}
							style={handleStyle(handle)}
							onPointerDown={(event) => event.stopPropagation()}
						/>
					)}
					resizeHandles={["sw", "se", "nw", "ne", "w", "e", "n", "s"]}
					onResizeStart={handleResizeStart}
					onResize={handleResize}
					onResizeStop={handleResizeStop}
					lockAspectRatio={false}
					axis="both"
					transformScale={1}
				>
					<div className="h-full w-full cursor-move rounded-md bg-black/10" onPointerDown={handleDragStart} />
				</ResizableBox>
			)}
		</div>
	);
}

export function FileSheet({
	trigger,
	username,
	file,
	post,
}: {
	trigger: ReactElement;
	username: string;
	file: string;
	post: ScrapeResponse;
}) {
	const [selectedTab, setSelectedTab] = useState<"view" | "crop">("view");
	const [cropRect, setCropRect] = useState<CropRect | null>(null);
	const [isFullImageCrop, setIsFullImageCrop] = useState(true);
	const [viewReloadKey, setViewReloadKey] = useState(0);
	const [resetSignal, setResetSignal] = useState(0);
	const cropFileMutation = useMutation(cropFile);
	const handleCropChange = useCallback((rect: CropRect | null, nextIsFullImageCrop: boolean) => {
		setCropRect(rect);
		setIsFullImageCrop(nextIsFullImageCrop);
	}, []);

	useEffect(() => {
		if (selectedTab !== "crop") {
			setCropRect(null);
			setIsFullImageCrop(true);
		}
	}, [selectedTab, file]);

	return (
		<Sheet>
			<SheetTrigger render={trigger} />
			<SheetContent side="bottom" className="p-1 pt-2 data-[side=bottom]:h-[90vh]">
				<Tabs value={selectedTab} onValueChange={setSelectedTab} className="flex h-full flex-col items-center">
					<TabsList>
						<TabsTrigger value="view">
							<ImageIcon />
							View
						</TabsTrigger>
						<TabsTrigger value="crop">
							<CropIcon />
							Crop
						</TabsTrigger>
					</TabsList>

					<div className="flex min-h-0 flex-1 flex-col items-center justify-center">
						<div className="flex max-h-full min-h-0 flex-col items-center">
							<TabsContent value="view" className="max-h-full w-auto">
								<FileDisplay
									key={viewReloadKey}
									file={file}
									post={post}
									username={username}
									className="h-full w-auto rounded-xl"
									cacheBuster={viewReloadKey}
								/>
							</TabsContent>
							<TabsContent
								value="crop"
								className="flex max-h-full min-h-0 w-auto flex-col items-center gap-2"
							>
								<ButtonGroup>
									<Button
										variant="outline"
										disabled={!cropFileMutation.isPending && (isFullImageCrop || cropRect === null)}
										onClick={() => {
											setResetSignal((prev) => prev + 1);
											setCropRect(null);
											setIsFullImageCrop(true);
										}}
									>
										<UndoIcon />
										Clear
									</Button>
									<Button
										variant="outline"
										disabled={!cropFileMutation.isPending && (isFullImageCrop || cropRect === null)}
										onClick={async () => {
											try {
												await cropFileMutation.mutateAsync({
													postType: post.postType,
													postOwner: post.postOwner,
													post: post.post,
													file,
													corner1: {
														x: cropRect?.x1,
														y: cropRect?.y1,
													},
													corner2: {
														x: cropRect?.x2,
														y: cropRect?.y2,
													},
												});
												setSelectedTab("view");
												const nextKey = Date.now();
												setViewReloadKey(nextKey);
												// Notify other components (e.g. parent FileDisplay instances) to bust cache for this file
												try {
													window.dispatchEvent(
														new CustomEvent("fileCropped", {
															detail: {
																username,
																postType: postTypeString(post.postType),
																postOwner: post.postOwner,
																file,
																cacheBuster: nextKey,
															},
														}),
													);
												} catch {}
											} catch (err) {
												toast.error((err as Error).message, {
													position: "top-center",
												});
											}
										}}
									>
										<CheckIcon />
										Done
									</Button>
								</ButtonGroup>
								{cropFileMutation.isPending && (
									<div className="w-full">
										<Progress value={null} />
									</div>
								)}
								<div className="flex min-h-0 flex-1 items-center justify-center">
									<CropPreview
										file={file}
										post={post}
										username={username}
										onCropChange={handleCropChange}
										cacheBuster={viewReloadKey}
										resetSignal={resetSignal}
									/>
								</div>
							</TabsContent>
						</div>
					</div>
				</Tabs>
			</SheetContent>
		</Sheet>
	);
}
