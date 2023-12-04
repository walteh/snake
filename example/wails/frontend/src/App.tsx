import { useEffect, useMemo, useState } from "react";
import { Commands, Inputs } from "../wailsjs/go/swails/WailsSnake";
import { swails } from "../wailsjs/go/models";
import Input from "./Input";
import Cards from "./Cards";

function App() {
	const [allCommands, setAllCommands] = useState<swails.WailsCommand[]>([]);

	const [allInputs, setAllInputs] = useState<swails.WailsInput[]>([]);

	useEffect(() => {
		Commands().then(setAllCommands);
		Inputs().then(setAllInputs);
	}, []);

	return (
		<div id="App">
			<div className="grid grid-cols-1 gap-4 sm:grid-cols-1 p-20">
				{allInputs?.map(
					(person) =>
						person.shared && (
							<div
								key={person.name}
								className="items-center space-x-3 rounded-lg bg-white  focus-within:ring-offset-2 "
							>
								<Input arg={person} />
							</div>
						)
				)}
			</div>
			<Cards commands={allCommands} />
		</div>
	);
}

export default App;
