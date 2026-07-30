import type { VariantProps } from "class-variance-authority";
import type { SVGProps } from "react";

import type { LatLng } from "@/buf/google/type/latlng_pb";
import { Button, type buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

function GoogleMaps(props: SVGProps<SVGSVGElement>) {
	return (
		<svg {...props} viewBox="0 0 92.3 132.3">
			<path fill="#1a73e8" d="M60.2 2.2C55.8.8 51 0 46.1 0 32 0 19.3 6.4 10.8 16.5l21.8 18.3L60.2 2.2z" />
			<path fill="#ea4335" d="M10.8 16.5C4.1 24.5 0 34.9 0 46.1c0 8.7 1.7 15.7 4.6 22l28-33.3-21.8-18.3z" />
			<path
				fill="#4285f4"
				d="M46.2 28.5c9.8 0 17.7 7.9 17.7 17.7 0 4.3-1.6 8.3-4.2 11.4 0 0 13.9-16.6 27.5-32.7-5.6-10.8-15.3-19-27-22.7L32.6 34.8c3.3-3.8 8.1-6.3 13.6-6.3"
			/>
			<path
				fill="#fbbc04"
				d="M46.2 63.8c-9.8 0-17.7-7.9-17.7-17.7 0-4.3 1.5-8.3 4.1-11.3l-28 33.3c4.8 10.6 12.8 19.2 21 29.9l34.1-40.5c-3.3 3.9-8.1 6.3-13.5 6.3"
			/>
			<path
				fill="#34a853"
				d="M59.1 109.2c15.4-24.1 33.3-35 33.3-63 0-7.7-1.9-14.9-5.2-21.3L25.6 98c2.6 3.4 5.3 7.3 7.9 11.3 9.4 14.5 6.8 23.1 12.8 23.1s3.4-8.7 12.8-23.2"
			/>
		</svg>
	);
}

export function GoogleMapsLink({
	coordinates,
	className,
	size = "sm",
}: {
	coordinates: LatLng;
	className?: string;
	size?: VariantProps<typeof buttonVariants>["size"];
}) {
	return (
		<Button
			variant="outline"
			size={size}
			className={cn("dark:bg-secondary dark:hover:bg-secondary/80", className)}
			nativeButton={false}
			render={
				<a
					href={`https://www.google.com/maps/search/?api=1&query=${coordinates.latitude},${coordinates.longitude}`}
					target="_blank"
					rel="noopener noreferrer"
					aria-label="Open location in Google Maps"
				/>
			}
		>
			<GoogleMaps />
		</Button>
	);
}

export function GooglePasswordMamagerIcon(props: SVGProps<SVGSVGElement>) {
	return (
		<svg {...props} viewBox="0 0 60 30">
			<path
				d="M20.9001 18.4125C19.7196 20.446 17.5223 21.8192 15.0035 21.8192C11.2459 21.8192 8.18363 18.7631 8.18363 15C8.18363 11.2369 11.24 8.18076 15.0035 8.18076C17.5223 8.18076 19.7196 9.55396 20.9001 11.5875H29.6076C28.059 4.9552 22.104 0 15.0035 0C6.73432 0 0.0020752 6.72575 0.0020752 15C0.0020752 23.2743 6.72848 30 15.0035 30C22.104 30 28.059 25.0448 29.6076 18.4125H20.9001Z"
				fill="#4285F4"
			/>
			<path d="M44.3238 10.915H29.3224V19.0958H44.3238V10.915Z" fill="#FBBC04" />
			<path
				d="M59.998 18.4126V27.277H54.5456V24.5482C54.5456 23.0406 53.3242 21.8193 51.8164 21.8193C50.3087 21.8193 49.0873 23.0406 49.0873 24.5482V27.277H43.6349V18.4126H59.998Z"
				fill="#34A853"
			/>
			<path d="M59.998 10.915H43.6349V19.0958H59.998V10.915Z" fill="#188038" />
			<path
				d="M29.4315 10.915H20.4552V10.9326C21.3084 12.072 21.8168 13.4803 21.8168 15.0054C21.8168 16.5305 21.3026 17.9388 20.4552 19.0783V19.0958H29.4315C29.7997 17.7927 30.0042 16.4254 30.0042 15.0054C30.0042 13.5855 29.7997 12.2181 29.4315 10.915Z"
				fill="#EA4335"
			/>
		</svg>
	);
}
