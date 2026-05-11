import { useCallback, useEffect, useRef, useState } from "react";

import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import {
	Carousel,
	CarouselContent,
	CarouselItem,
	CarouselNext,
	CarouselPrevious,
	type CarouselApi,
} from "@/components/ui/carousel";
import { GoogleMapsLink } from "@/components/ui/svgs/google-maps";

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
			withCoordinates
		/>
	);
}

export function FileDisplay({
	username,
	file,
	post: { postType, postOwner, coordinates },
	onMediaLoad,
	withCoordinates,
}: {
	username: string;
	file: string;
	post: ScrapeResponse;
	onMediaLoad?: () => void;
	withCoordinates?: boolean;
}) {
	const url = `/api/storage/${username}/${postTypeString(postType)}/${postOwner}/${file}`;
	if (/\.(jpg)|(jpeg)|(webp)|(heic)$/.test(file)) {
		const imgResult = <img src={url} onLoad={onMediaLoad} loading="lazy" className="h-auto w-full" />;
		return withCoordinates && postType === PostType.VSCO && coordinates ? (
			<div className="relative inline-block w-full">
				{imgResult}
				<GoogleMapsLink coordinates={coordinates} className="absolute top-2 left-2 z-10" size="icon" />
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
				className="h-auto w-full"
				loop
				controls
				muted
			/>
		);
	} else {
		return <a href={url}>{url}</a>;
	}
}
