:- encoding(utf8).

%%%%%%%%%
%%%%%%%%%
:- use_module(library(assoc)).

dic_empty(D) :- empty_assoc(D).

dic_get(D, K, V) :- get_assoc(K, D, V).

dic_replace(D0, K, V, D) :- put_assoc(K, D0, V, D).

%%%%%%%%%
%%%%%%%%%
interpret_nv(Program, AcIn, AcOut):-
	dic_empty(DicIn),
	interpret_naive(Program, Program, DicIn, _DicOut, AcIn, AcOut).

%%%%%%%%%
interpret_naive([], _Program, _, _, Acum, Acum).

interpret_naive([LI|Is], Program, DictIn, DictOut, AIn, AOut):-
	remove_label(LI, I),
	execute(I, Is, Program, NProgram, DictIn, DictMid, AIn, AMid),
	interpret_naive(NProgram, Program, DictMid, DictOut, AMid, AOut).

%%%%%%%%%
execute(load(NumOrLabel), Is, _Program, Is, D, D, _, Value):-
	eval(NumOrLabel, D, Value).

execute(add(NumOrLabel), Is, _Program, Is, D, D, AcIn, AcOut):-
	eval(NumOrLabel, D, Value),
	AcOut is AcIn + Value.

execute(sub(NumOrLabel), Is, _Program, Is, D, D, AcIn, AcOut):-
	eval(NumOrLabel, D, Value),
	AcOut is AcIn - Value.

%%
execute(sto(Label), Is, _Program, Is, DIn, DOut, Acum, Acum):-
	dic_replace(DIn, Label, Acum, DOut).

%%
execute(jmp(Label), _, Program, Is, D, D, A, A):-
	Is = [Label:_|_], append(_, Is, Program).

execute(jez(Label), NIs, Program, Is, D, D, A, A):-
	(A = 0 -> Is = [Label:_|_], append(_, Is, Program) ; Is = NIs).

execute(jnez(Label), NIs, Program, Is, D, D, A, A):-
	(A \== 0 -> Is = [Label:_|_], append(_, Is, Program) ; Is = NIs).

%%
execute(nop, Is, _Program, Is, D, D, A, A).

%%%%%%%%%
remove_label(_:I, I):- !.
remove_label(I, I).

%%%%%%%%%
eval(Number, _Dict, Number):-
	number(Number).

eval(Label, Dict, Number):-
	atom(Label),
	dic_get(Dict, Label, Number).


%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%
program(fibo, [
jnez(calculate),

load(0), sto(curr), jmp(end),

calculate:sto(ind),
		load(0),
		sto(prev),
		load(1),
		sto(curr),

start_loop:load(ind),
		sub(1),
		sto(ind), 
		jez(end),
		load(curr),
		sto(inter),
		add(prev), 
		sto(curr),
		load(inter),
		sto(prev), 
		jmp(start_loop), 

end:load(curr)
]).

%%
main :- program(fibo, Prg),
	interpret_nv(Prg, 400, A),
	writeln(A).
