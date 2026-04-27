import { MatchRoute } from "@tanstack/react-router";
import { SidebarIcon } from "lucide-react";

import { RakerLogo } from "@/components/logo";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";
import { ThemeToggle } from "@/components/ui/theme-toggle";

export default function Header({ toggleMenu }: { toggleMenu: () => void }) {
	return (
		<header className="sticky top-0 z-50 flex w-full items-center border-b bg-background backdrop-blur-lg">
			<div className="grid h-(--header-height) w-full grid-cols-[1fr_auto_1fr] items-center px-4">
				<div className="justify-self-start">
					<Button className="h-8 w-8" variant="outline" size="icon" onClick={toggleMenu}>
						<SidebarIcon />
					</Button>
				</div>
				<div className="justify-self-center">
					<div className="flex flex-row items-center *:mx-1">
						<RakerLogo withVersion />
						<MatchRoute to="/instagram">
							<Separator orientation="vertical" className="" />
							<InstagramIcon className="h-4" />
							<Label>Post</Label>
						</MatchRoute>
						<MatchRoute to="/highlight">
							<Separator orientation="vertical" className="" />
							<InstagramIcon className="h-4" />
							<Label>Highlight</Label>
						</MatchRoute>
						<MatchRoute to="/story">
							<Separator orientation="vertical" className="" />
							<InstagramIcon className="h-4" />
							<Label>Story</Label>
						</MatchRoute>
						<MatchRoute to="/tiktok">
							<Separator orientation="vertical" className="" />
							<TikTokIcon className="h-4" />
							<Label>Post</Label>
						</MatchRoute>
						<MatchRoute to="/vsco">
							<Separator orientation="vertical" className="" />
							<VSCOIcon className="h-4" />
							<Label>Post</Label>
						</MatchRoute>
					</div>
				</div>
				<div className="justify-self-end">
					<ThemeToggle />
				</div>
			</div>
		</header>
	);
}
