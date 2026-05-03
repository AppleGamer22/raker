import { Progress as ProgressPrimitive } from "@base-ui/react/progress";

import { cn } from "@/lib/utils";

function Progress({ className, children, value, ...props }: ProgressPrimitive.Root.Props) {
	const isIndeterminate = value == null;

	return (
		<ProgressPrimitive.Root
			value={value}
			data-indeterminate={isIndeterminate ? "true" : "false"}
			data-slot="progress"
			className={cn("flex flex-wrap gap-3", className)}
			{...props}
		>
			{children}
			<ProgressTrack>
				{isIndeterminate ? (
					<span
						aria-hidden="true"
						className="animate-progress-indeterminate absolute inset-y-0 left-0 w-2/5 rounded-full bg-primary"
					/>
				) : (
					<ProgressIndicator />
				)}
			</ProgressTrack>
		</ProgressPrimitive.Root>
	);
}

function ProgressTrack({ className, ...props }: ProgressPrimitive.Track.Props) {
	return (
		<ProgressPrimitive.Track
			className={cn("relative flex h-1 w-full items-center overflow-x-hidden rounded-full bg-muted", className)}
			data-slot="progress-track"
			{...props}
		/>
	);
}

function ProgressIndicator({ className, ...props }: ProgressPrimitive.Indicator.Props) {
	return (
		<ProgressPrimitive.Indicator
			data-slot="progress-indicator"
			className={cn("h-full bg-primary transition-all", className)}
			{...props}
		/>
	);
}

function ProgressLabel({ className, ...props }: ProgressPrimitive.Label.Props) {
	return (
		<ProgressPrimitive.Label
			className={cn("text-sm font-medium", className)}
			data-slot="progress-label"
			{...props}
		/>
	);
}

function ProgressValue({ className, ...props }: ProgressPrimitive.Value.Props) {
	return (
		<ProgressPrimitive.Value
			className={cn("ml-auto text-sm text-muted-foreground tabular-nums", className)}
			data-slot="progress-value"
			{...props}
		/>
	);
}

export { Progress, ProgressTrack, ProgressIndicator, ProgressLabel, ProgressValue };
