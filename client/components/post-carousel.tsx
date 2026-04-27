import { useEffect, useState } from "react";

import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import {
	Carousel,
	CarouselContent,
	CarouselItem,
	CarouselNext,
	CarouselPrevious,
	type CarouselApi,
} from "@/components/ui/carousel";

function postTypeString(type: PostType) {
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

export function PostCarousel({
	username,
	post: { postType, postOwner, post, files },
}: {
	username: string;
	post: ScrapeResponse;
}) {
	const [api, setApi] = useState<CarouselApi>();
	const [current, setCurrent] = useState(0);
	const [count, setCount] = useState(0);

	useEffect(() => {
		if (!api) {
			return;
		}

		setCount(api.scrollSnapList().length);
		setCurrent(api.selectedScrollSnap() + 1);

		api.on("select", () => {
			setCurrent(api.selectedScrollSnap() + 1);
		});
	}, [api]);

	return (
		<>
			<Carousel setApi={setApi} className="w-full">
				<CarouselContent>
					{files.map((file, i) => (
						<CarouselItem key={`file-${postType}-${postOwner}-${post}-${i}`}>
							<img src={`/api/storage/${username}/${postTypeString(postType)}/${postOwner}/${file}`} />
						</CarouselItem>
					))}
				</CarouselContent>
				<CarouselPrevious className="top-1/2 left-2 z-10 -translate-y-1/2" />
				<CarouselNext className="top-1/2 right-2 z-10 -translate-y-1/2" />
			</Carousel>
			<div className="py-2 text-center text-sm text-muted-foreground">
				Slide {current} of {count}
			</div>
		</>
	);
}
