-- Remove sample mentors
DELETE FROM mentors WHERE email IN (
    't.salakhov@innopolis.university',
    'i.abdulkhakov@innopolis.university',
    'd.minnakhmetova@innopolis.university',
    'v.zhidkov@innopolis.university'
);
