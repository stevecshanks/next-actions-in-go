interface ActionArgs {
  id: string;
  name: string;
  dueBy?: Date;
}

export class Action {
  id: string;
  name: string;
  dueBy?: Date;

  constructor(args: ActionArgs) {
    const { id, name, dueBy } = args;

    this.id = id;
    this.name = name;
    this.dueBy = dueBy;
  }
}
