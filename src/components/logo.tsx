import { Badge } from "@/components/ui/badge";

export function RakerLogo({ withVersion }: { withVersion?: boolean }) {
	return (
		<div className="flex flex-row items-center">
			<img alt="Raker Logo" src="/raker.svg" className="w-6" />
			{withVersion && (
				<sup className="ml-2">
					<Badge>&alpha;</Badge>
				</sup>
			)}
		</div>
	);
}
