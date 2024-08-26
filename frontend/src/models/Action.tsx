export interface ActionArgs {
  id: string;
  name: string;
  url: string;
  imageUrl: string | null;
  projectName: string;
  dueBy?: Date;
}

const TWENTY_FOUR_HOURS = 24 * 60 * 60 * 1000;

export class Action {
  id: string;
  name: string;
  url: string;
  imageUrl: string | null;
  projectName: string;
  dueBy?: Date;

  constructor(args: ActionArgs) {
    const { id, name, url, imageUrl, projectName, dueBy } = args;

    this.id = id;
    this.name = name;
    this.url = url;
    this.imageUrl = imageUrl;
    this.projectName = projectName;
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

  isDueThisWeek(): boolean {
    if (!this.dueBy) {
      return false;
    }
    return this.dueBy.getTime() <= endOfCurrentWeek().getTime();
  }
}

const endOfCurrentWeek = (): Date => {
  const today = new Date();
  const mondayDate = today.getDate() - today.getDay() + 1;
  const sundayDate = new Date(today.setDate(mondayDate + 6));
  sundayDate.setHours(23, 59, 59, 999);
  return sundayDate;
};
