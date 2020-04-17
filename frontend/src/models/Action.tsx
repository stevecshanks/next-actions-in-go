interface ActionArgs {
  id: string;
  name: string;
  url: string;
  imageUrl: string;
  dueBy?: Date;
}

const TWENTY_FOUR_HOURS = 24 * 60 * 60 * 1000;

export class Action {
  id: string;
  name: string;
  url: string;
  imageUrl: string;
  dueBy?: Date;

  constructor(args: ActionArgs) {
    const { id, name, url, imageUrl, dueBy } = args;

    this.id = id;
    this.name = name;
    this.url = url;
    this.imageUrl = imageUrl;
    this.dueBy = dueBy;
  }

  isOverdue(): boolean {
    if (!this.dueBy) {
      return false;
    }
    return this.dueBy.getTime() < Date.now();
  }

  isDueSoon(): boolean {
    if (!this.dueBy) {
      return false;
    }
    return this.dueBy.getTime() < Date.now() + TWENTY_FOUR_HOURS;
  }
}
