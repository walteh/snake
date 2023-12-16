import { useCallback, useEffect, useRef } from "react";
import { Terminal } from "xterm";
import { swails } from "../wailsjs/go/models";
import { EventsOn } from "../wailsjs/runtime";
import "xterm/css/xterm.css";

export const useXTerm = () => {
	// The element to create the terminal within. This element must be visible (have dimensions) when open is called as several DOM- based measurements need to be performed when this function is called.
	const terminalRef = useRef<HTMLDivElement>(null);
	const term = useRef(new Terminal({})); // Create a ref to hold the Terminal instance

	// Function to write to the terminal
	const writeToTerminal = useCallback((text: string) => {
		term.current.write(text);
	}, []);

	// Component to render the terminal
	const XTermComponent = () => {
		// Initialize the terminal
		useEffect(() => {
			if (!terminalRef.current) {
				return;
			}
			// console.log("opening terminal", { ...terminalRef.current });
			term.current.open(terminalRef.current);

			// Optional cleanup
			return () => {
				term.current.dispose();
			};
		}, [terminalRef.current]);

		return (
			<div ref={terminalRef} style={{ height: "100%", width: "100%" }} />
		);
	};
	// Assume you have this situation: Bob orders an egg mcmuffin today and gets an offer for a half off one tomorrow. Currently
	return [XTermComponent, writeToTerminal] as const;
};

export const XTermWindow = ({ writer }: { writer: swails.WailsWriter }) => {
	const [XTermComponent, writeToTerminal] = useXTerm();

	useEffect(() => {
		const handleWriteToTerminal = (arg: any) => {
			console.log("data: ", arg);
			// The data is in e.detail
			writeToTerminal(arg);
		};

		console.log("opening writer: ", writer.name);

		const clean = EventsOn(writer.name, handleWriteToTerminal);

		// listen for data events from the writer by checking for events withthe name of the writer
		// const clean = emitter.on(writer.name, handleWriteToTerminal);

		console.log("event listening");

		// Cleanup
		return () => {
			console.log("closing writer: ", writer.name);
			clean();
			// clean.removeListener(writer.name, handleWriteToTerminal);
		};
	}, []);

	return (
		<div>
			<XTermComponent />
		</div>
	);
};

export default XTermWindow;
