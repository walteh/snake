import { useEffect, useMemo, useState } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import {
	InputsFor,
	Run,
	Commands,
	Inputs,
	UpdateInput,
} from "../wailsjs/go/swails/WailsSnake";
import { swails } from "../wailsjs/go/models";
import Input from "./Input";

function App() {
	const [resultText, setResultText] = useState<swails.WailsHTMLResponse>(
		new swails.WailsHTMLResponse()
	);
	const [name, setName] = useState("");
	const updateName = (e: any) => setName(e.target.value);
	const updateResultText = (result: swails.WailsHTMLResponse) =>
		setResultText(result);

	const [allInputs, setAllInputs] =
		useState<Map<string, swails.WailsInput>>();

	const [allCommands, setAllCommands] =
		useState<Map<string, swails.WailsCommand>>();

	useEffect(() => {
		Commands().then((result) => {
			console.log("Commands: ", { result });
			const obj: Map<string, swails.WailsCommand> = new Map();
			for (let i = 0; i < result.length; i++) {
				obj.set(result[i].name, result[i]);
			}
			setAllCommands(obj);
		});

		Inputs().then((result) => {
			const obj: Map<string, swails.WailsInput> = new Map();
			console.log("Inputs: ", { result });
			for (let i = 0; i < result.length; i++) {
				obj.set(result[i].name, result[i]);
			}
			setAllInputs(obj);
		});
	}, []);

	function greet(cmd: swails.WailsCommand) {
		Run(cmd).then((result) => {
			console.log(result);
			updateResultText(result);
		});
	}

	const response = useMemo(() => {
		return (
			<div
				id="result"
				className="result"
				dangerouslySetInnerHTML={{
					__html: resultText.html,
				}}
			></div>
		);
	}, [resultText]);

	return (
		<div id="App">
			{allCommands && (
				<div>
					{Array.from(allCommands.keys()).map((key) => (
						<div key={key}>
							{key} : {allCommands.get(key)?.description}{" "}
							<button
								className="btn"
								onClick={() => {
									if (allCommands.get(key) !== undefined) {
										greet(
											allCommands.get(
												key
											) as swails.WailsCommand
										);
									}
								}}
							>
								Greet
							</button>
						</div>
					))}
				</div>
			)}
			{response}

			<div id="input" className="input-box">
				{allInputs && (
					<div>
						{Array.from(allInputs.keys()).map((key) => (
							<div key={key}>
								<Input
									arg={
										allInputs.get(key) as swails.WailsInput
									}
								/>
							</div>
						))}
					</div>
				)}
			</div>
		</div>
	);
}

export default App;
